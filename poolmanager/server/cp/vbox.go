package cp

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/google/uuid"
)

const vboxNamePrefix = "fn-vagrant"

var whichVBox *exec.Cmd

func init() {
	whichVBox = exec.Command("which", "vagrant")
}

type VirtualBoxCP struct{}

func NewVirtualBoxCP() (*VirtualBoxCP, error) {
	if err := whichVBox.Run(); err != nil {
		return nil, err
	}
	return &VirtualBoxCP{}, nil
}

func (v *VirtualBoxCP) provision() error {
	//set up dir
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	node := newNodeName()
	nodeDir, err := ioutil.TempDir(wd, node)
	if err != nil {
		return err
	}
	//copy vagrant file into there
	vagrantFile := fmt.Sprintf("%s/%s", wd, "Vagrantfile")
	newVagrantFile := fmt.Sprintf("%s/%s", nodeDir, "Vagrantfile")
	err = copyFile(vagrantFile, newVagrantFile)
	if err != nil {
		return err
	}

	err = os.Chdir(nodeDir)
	if err != nil {
		return err
	}
	vboxProvision := exec.Command("vagrant", "up")
	err = vboxProvision.Run()
	if err != nil {
		return err
	}
	return nil
}

func (v *VirtualBoxCP) GetLBGRunners(lgbId string) ([]*Runner, error) {
	return nil, errors.New("Not done")

}

func (v *VirtualBoxCP) ProvisionRunners(lgbId string, n int) (int, error) {
	return -1, errors.New("Not done")
}

func (v *VirtualBoxCP) RemoveRunner(lbgId string, id string) error {
	return errors.New("Not done")
}

func newNodeName() string {
	id := uuid.New()
	return fmt.Sprintf("%s-%s", vboxNamePrefix, id.String())
}

func copyFile(src string, dst string) error {
	// Open the source file for reading
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	// Open the destination file for writing
	d, err := os.Create(dst)
	if err != nil {
		return err
	}

	// Copy the contents of the source file into the destination file
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}

	// Return any errors that result from closing the destination file
	// Will return nil if no errors occurred
	return d.Close()
}
