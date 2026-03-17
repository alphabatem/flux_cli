package services

import (
	stdcontext "context"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/alphabatem/flux_cli/dto"
	pb "github.com/alphabatem/flux_cli/internal/yellowstonepb"
	ctxpkg "github.com/alphabatem/flux_cli/pkg/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
)

const YELLOWSTONE_SVC = "yellowstone_svc"

type YellowstoneService struct {
	ctxpkg.DefaultService

	mu     sync.Mutex
	client pb.GeyserClient
	conn   *grpc.ClientConn
	cfg    dto.FluxRPCConfig
}

func (s *YellowstoneService) Id() string {
	return YELLOWSTONE_SVC
}

func ResolveYellowstoneURL(cfg *dto.FluxRPCConfig) string {
	switch cfg.Region {
	case "eu":
		return "https://yellowstone.eu.fluxrpc.com"
	case "us":
		fallthrough
	default:
		return "https://yellowstone.us.fluxrpc.com"
	}
}

func (s *YellowstoneService) Configure(ctx *ctxpkg.Context) error {
	cfg := ctx.Service(CONFIG_SVC).(*ConfigService).Config()
	s.cfg = cfg.FluxRPC
	return s.DefaultService.Configure(ctx)
}

func (s *YellowstoneService) Shutdown() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn != nil {
		_ = s.conn.Close()
		s.conn = nil
		s.client = nil
	}
}

func (s *YellowstoneService) WatchAccounts(ctx stdcontext.Context, accounts []string, commitment pb.CommitmentLevel, onUpdate func(*pb.SubscribeUpdate) error) error {
	commitmentCopy := commitment
	req := &pb.SubscribeRequest{
		Accounts: map[string]*pb.SubscribeRequestFilterAccounts{
			"accounts": {Account: accounts},
		},
		Commitment: &commitmentCopy,
	}

	return s.watch(ctx, req, onUpdate)
}

func (s *YellowstoneService) WatchProgramOwners(ctx stdcontext.Context, owners []string, commitment pb.CommitmentLevel, onUpdate func(*pb.SubscribeUpdate) error) error {
	commitmentCopy := commitment
	req := &pb.SubscribeRequest{
		Accounts: map[string]*pb.SubscribeRequestFilterAccounts{
			"owners": {Owner: owners},
		},
		Commitment: &commitmentCopy,
	}

	return s.watch(ctx, req, onUpdate)
}

func (s *YellowstoneService) WatchSlots(ctx stdcontext.Context, commitment pb.CommitmentLevel, interslotUpdates bool, onUpdate func(*pb.SubscribeUpdate) error) error {
	interslot := interslotUpdates
	commitmentCopy := commitment
	req := &pb.SubscribeRequest{
		Slots: map[string]*pb.SubscribeRequestFilterSlots{
			"slots": {InterslotUpdates: &interslot},
		},
		Commitment: &commitmentCopy,
	}

	return s.watch(ctx, req, onUpdate)
}

func (s *YellowstoneService) WatchTransactions(
	ctx stdcontext.Context,
	accountInclude []string,
	accountExclude []string,
	accountRequired []string,
	includeVotes bool,
	includeFailed bool,
	commitment pb.CommitmentLevel,
	onUpdate func(*pb.SubscribeUpdate) error,
) error {
	commitmentCopy := commitment
	vote := includeVotes
	failed := includeFailed
	req := &pb.SubscribeRequest{
		Transactions: map[string]*pb.SubscribeRequestFilterTransactions{
			"transactions": {
				Vote:            &vote,
				Failed:          &failed,
				AccountInclude:  accountInclude,
				AccountExclude:  accountExclude,
				AccountRequired: accountRequired,
			},
		},
		Commitment: &commitmentCopy,
	}

	return s.watch(ctx, req, onUpdate)
}

func (s *YellowstoneService) WatchTransactionSignature(
	ctx stdcontext.Context,
	signature string,
	commitment pb.CommitmentLevel,
	onUpdate func(*pb.SubscribeUpdate) error,
) error {
	commitmentCopy := commitment
	sig := strings.TrimSpace(signature)
	req := &pb.SubscribeRequest{
		Transactions: map[string]*pb.SubscribeRequestFilterTransactions{
			"signature": {
				Signature: &sig,
			},
		},
		Commitment: &commitmentCopy,
	}

	return s.watch(ctx, req, onUpdate)
}

func (s *YellowstoneService) watch(ctx stdcontext.Context, req *pb.SubscribeRequest, onUpdate func(*pb.SubscribeUpdate) error) error {
	backoff := 300 * time.Millisecond

	for {
		if ctx.Err() != nil {
			return nil
		}

		err := s.streamOnce(ctx, req, onUpdate)
		if err == nil || ctx.Err() != nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(backoff):
		}

		if backoff < 5*time.Second {
			backoff *= 2
		}

		s.disconnect()
	}
}

func (s *YellowstoneService) streamOnce(ctx stdcontext.Context, req *pb.SubscribeRequest, onUpdate func(*pb.SubscribeUpdate) error) error {
	client, err := s.ensureClient()
	if err != nil {
		return err
	}

	streamCtx := ctx
	if s.cfg.APIKey != "" {
		streamCtx = metadata.NewOutgoingContext(ctx, metadata.Pairs("x-token", s.cfg.APIKey))
	}

	stream, err := client.Subscribe(streamCtx)
	if err != nil {
		return err
	}
	defer stream.CloseSend()

	if err := stream.Send(req); err != nil {
		return err
	}

	for {
		update, recvErr := stream.Recv()
		if recvErr != nil {
			if errors.Is(recvErr, io.EOF) {
				return nil
			}
			return recvErr
		}

		// Keepalive/ping responses do not include user data.
		if update == nil || len(update.Filters) == 0 {
			continue
		}

		if err := onUpdate(update); err != nil {
			return err
		}
	}
}

func (s *YellowstoneService) ensureClient() (pb.GeyserClient, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.client != nil {
		return s.client, nil
	}

	if strings.TrimSpace(s.cfg.APIKey) == "" {
		return nil, errors.New("FluxRPC API key not configured. Run: flux config set fluxrpc.api_key <key>")
	}

	target := ResolveYellowstoneURL(&s.cfg)
	conn, err := connectGRPC(target)
	if err != nil {
		return nil, err
	}

	s.conn = conn
	s.client = pb.NewGeyserClient(conn)
	return s.client, nil
}

func (s *YellowstoneService) disconnect() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn != nil {
		_ = s.conn.Close()
	}
	s.conn = nil
	s.client = nil
}

func connectGRPC(rawURL string) (*grpc.ClientConn, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("parsing yellowstone URL: %w", err)
	}

	httpMode := u.Scheme == "http"

	port := u.Port()
	if port == "" {
		if httpMode {
			port = "80"
		} else {
			port = "443"
		}
	}
	host := u.Hostname()
	if host == "" {
		return nil, errors.New("invalid yellowstone URL: expected http(s)://host[:port]")
	}
	address := host + ":" + port

	opts := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	}

	if httpMode {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		pool, _ := x509.SystemCertPool()
		creds := credentials.NewClientTLSFromCert(pool, "")
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("connecting to yellowstone: %w", err)
	}
	return conn, nil
}
