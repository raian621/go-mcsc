package api

import (
	"encoding/json"
	"io"
	"net/http"
)

type ServerController struct {
	msi MinecraftServerInterface
}

// GetAllowlist implements ServerInterface.
func (s *ServerController) GetAllowlist(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// GetAvailableVersions implements ServerInterface.
func (s *ServerController) GetAvailableVersions(w http.ResponseWriter, r *http.Request) {
	versions := s.msi.Versions()

	if versions == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(versions); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GetBannedIps implements ServerInterface.
func (s *ServerController) GetBannedIps(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// GetBannedPlayers implements ServerInterface.
func (s *ServerController) GetBannedPlayers(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// GetOps implements ServerInterface.
func (s *ServerController) GetOps(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PostAllowlistAdd implements ServerInterface.
func (s *ServerController) PostAllowlistAdd(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PostAllowlistRemove implements ServerInterface.
func (s *ServerController) PostAllowlistRemove(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PostBan implements ServerInterface.
func (s *ServerController) PostBan(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PostBanIp implements ServerInterface.
func (s *ServerController) PostBanIp(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PostDeop implements ServerInterface.
func (s *ServerController) PostDeop(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PostOp implements ServerInterface.
func (s *ServerController) PostOp(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PostPardon implements ServerInterface.
func (s *ServerController) PostPardon(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PostPardonIp implements ServerInterface.
func (s *ServerController) PostPardonIp(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PostRestart implements ServerInterface.
func (s *ServerController) PostRestart(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PostSetVersion implements ServerInterface.
func (s *ServerController) PostSetVersion(w http.ResponseWriter, r *http.Request, params PostSetVersionParams) {
	panic("unimplemented")
}

// PostStart implements ServerInterface.
func (s *ServerController) PostStart(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PostStop implements ServerInterface.
func (s *ServerController) PostStop(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PutAllowlist implements ServerInterface.
func (s *ServerController) PutAllowlist(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PutArgs implements ServerInterface.
func (s *ServerController) PutArgs(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PutBannedIps implements ServerInterface.
func (s *ServerController) PutBannedIps(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PutBannedPlayers implements ServerInterface.
func (s *ServerController) PutBannedPlayers(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PutOps implements ServerInterface.
func (s *ServerController) PutOps(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PutProperties implements ServerInterface.
func (s *ServerController) PutProperties(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

type MinecraftServerConfig struct {
	Version string `json:"version,omitempty"`
}

type MinecraftServerInterface interface {
	// allowlist methods

	Allowlist() *Allowlist
	AllowPlayer(p *PlayerInfo) error
	DisallowPlayer(p *PlayerInfo) error

	// ban and unban methods

	BannedIPs() *BannedIPList
	BannedPlayers() *BannedPlayerList
	BanIP(ip *BannedIP) error
	BanPlayer(p *BannedPlayer) error
	PardonIP(ip string) error
	PardonPlayer(p *PlayerInfo) error

	// config files initialization methods

	CreateAllowlist()
	CreateArgs()
	CreateBannedIPs()
	CreateBannedPlayers()
	CreateConfig()
	CreateProperties()
	CreateOperators()
	CreateVersions()

	// server operator methods

	Deop(p *PlayerInfo) error
	Op(op *ServerOperator) error
	Ops() *ServerOperatorList

	// config files loading methods

	LoadAllowlist(file io.Reader) error
	LoadArgs(file io.Reader) error
	LoadBannedIPs(file io.Reader) error
	LoadBannedPlayers(file io.Reader) error
	LoadConfig(file io.Reader) error
	LoadConfigs() error
	LoadOperators(file io.Reader) error
	LoadProperties(file io.Reader) error
	LoadVersions(file io.Reader) error

	// miscellaneous get config methods:

	Args() *ServerArguments
	Config() *MinecraftServerConfig
	Properties() *ServerProperties
	Versions() *[]string

	// save config files methods

	SaveAllowlist(file io.Writer) error
	SaveArgs(file io.Writer) error
	SaveBannedIPs(file io.Writer) error
	SaveBannedPlayers(file io.Writer) error
	SaveConfig(file io.Writer) error
	SaveOperators(file io.Writer) error
	SaveProperties(file io.Writer) error

	// update server configs methods (doesn't save to disk)

	SetAllowlist(a *Allowlist)
	SetArgs(args *ServerArguments)
	SetBannedPlayers(bp *BannedPlayerList)
	SetVersion(version string) error
	SetBannedIPs(bp *BannedIPList)
	SetOperators(ops *ServerOperatorList)
	SetProperties(props *ServerProperties)

	// mc server process management methods

	Start() error
	Stop() error
	Restart() error
}

var _ ServerInterface = (*ServerController)(nil)

func NewServerController(msi MinecraftServerInterface) ServerInterface {
	return &ServerController{msi}
}
