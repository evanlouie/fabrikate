package ssh

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/microsoft/fabrikate/internal/strutil"
)

// PlainIdentity encapsulates the output of `ssh-add`
type PlainIdentity struct {
	Path    string
	Comment string
}

// Fingerprint encapsulates the output of `ssh-add -l`
type Fingerprint struct {
	BitLength      int    // 256, 4096, etc..
	Fingerprint    string // SHA256:MgTreyI8MtnEN9Yh1KrNdxHFM2wqOwnWNIGiRDTDCW8
	Comment        string // foo@bar.com
	EncryptionType string // RSA, ED25519, etc...
}

// InitializeIdentities of the host SSH agent by shelling out to `ssh-add`.
func InitializeIdentities() (identities []PlainIdentity, err error) {
	cmd := exec.Command("ssh-add")
	// NOTE `ssh-add` outputs everything to stderr
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf(`running "%s": %s: %w`, cmd, string(out), err)
	}

	// prepare the regex
	var (
		Path           = "Path"
		Comment        = "Comment"
		identityRgxStr = fmt.Sprintf(`(?i)Identity added: (?P<%s>\S+) \((?P<%s>\S+)\)`, Path, Comment)
		identityRgx    = regexp.MustCompile(identityRgxStr)
	)
	for _, line := range strutil.SplitLines(string(out)) {
		if identityRgx.MatchString(line) {
			var identity PlainIdentity
			for idx, match := range identityRgx.FindStringSubmatch(line) {
				switch identityRgx.SubexpNames()[idx] {
				case Path:
					identity.Path = match
				case Comment:
					identity.Comment = match
				}
			}
			identities = append(identities, identity)
		}
	}

	return identities, err
}

// Fingerprints returns the output of `ssh-add -l`
func Fingerprints() (prints []Fingerprint, err error) {
	cmd := exec.Command("ssh-add", "-l")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf(`running "%s": %s: %w`, cmd, out, err)
	}

	var (
		BitLength  = "BitLength"
		Print      = "Print"
		Comment    = "Comment"
		Encryption = "Encryption"
		rgxString  = fmt.Sprintf(`(?i)(?P<%s>\d+) (?P<%s>\S+)(\s(?P<%s>\S+))? \((?P<%s>\S+)\)`, BitLength, Print, Comment, Encryption)
		rgx        = regexp.MustCompile(rgxString)
	)
	for _, line := range strutil.SplitLines(string(out)) {
		if rgx.MatchString(line) {
			var print Fingerprint
			for idx, match := range rgx.FindStringSubmatch(line) {
				switch rgx.SubexpNames()[idx] {
				case BitLength:
					matchAsInt, err := strconv.Atoi(match)
					if err != nil {
						return nil, fmt.Errorf(`converting ssh key bit length %s to integer: %w`, match, err)
					}
					print.BitLength = matchAsInt
				case Print:
					print.Fingerprint = match
				case Comment:
					print.Comment = match
				case Encryption:
					print.EncryptionType = match
				}
			}
			prints = append(prints, print)
		}
	}

	return prints, err
}
