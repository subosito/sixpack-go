package sixpack

import (
	"testing"
)

// return an alternative for participate
func TestAlternativeForParticipate(t *testing.T) {
	session, err := NewSession(Options{
		ClientID: []byte("mike"),
	})

	if err != nil {
		t.Error(err)
	}

	alternatives := []string{"trolled", "not-trolled"}

	r, err := session.Participate("show-bieber", alternatives, "")
	if err != nil {
		t.Error(err)
	}

	exist := false

	for _, v := range alternatives {
		if v == r.Alternative.Name {
			exist = true
		}
	}

	if !exist {
		t.Errorf("Alternative name: %q is not a member of %q", r.Alternative.Name, alternatives)
	}
}

// return the correct alternative for participate with force
func TestAlternativeForParticipateWithForce(t *testing.T) {
	session, err := NewSession(Options{
		ClientID: []byte("mike"),
	})

	if err != nil {
		t.Error(err)
	}

	alternatives := []string{"trolled", "not-trolled"}
	force := "trolled"

	r, err := session.Participate("show-bieber", alternatives, force)
	if err != nil {
		t.Error(err)
	}

	if force != r.Alternative.Name {
		t.Errorf("Force %q != %q", force, r.Alternative.Name)
	}
}

// Allow ip and user agetn to be passed to a session
func TestPassingOptions(t *testing.T) {
	session, err := NewSession(Options{
		IP:        "8.8.8.8",
		UserAgent: "FirChromari",
	})

	if err != nil {
		t.Error(err)
	}

	if tip := "8.8.8.8"; session.Options.IP != tip {
		t.Errorf("session.IP => %s != %s\n", session.Options.IP, tip)
	}
}

// Auto generate a clientID
func TestAutoGenerateClientID(t *testing.T) {
	session, err := NewSession(Options{})
	if err != nil {
		t.Error(err)
	}

	if length := 36; len(session.Options.ClientID) != length {
		t.Errorf("length session.Options.ClientID != %d", length)
	}
}

// Return ok for convert
func TestOkForConvert(t *testing.T) {
	session, err := NewSession(Options{
		ClientID: []byte("mike"),
	})

	if err != nil {
		t.Error(err)
	}

	_, err = session.Participate("show-bieber", []string{"trolled", "not-trolled"}, "")
	if err != nil {
		t.Error(err)
	}

	rc, err := session.Convert("show-bieber")
	if err != nil {
		t.Error(err)
	}

	if status := "ok"; rc.Status != status {
		t.Errorf("Status: %q != %q", rc.Status, status)
	}
}

// Return ok for multiple converts
func TestOkForMultipleConverts(t *testing.T) {
	session, err := NewSession(Options{
		ClientID: []byte("mike"),
	})

	if err != nil {
		t.Error(err)
	}

	_, err = session.Participate("show-bieber", []string{"trolled", "not-trolled"}, "")
	if err != nil {
		t.Error(err)
	}

	rc, err := session.Convert("show-bieber")
	if err != nil {
		t.Error(err)
	}

	if status := "ok"; rc.Status != status {
		t.Errorf("Status: %q != %q", rc.Status, status)
	}

	rcs, err := session.Convert("show-bieber")
	if err != nil {
		t.Error(err)
	}

	if status := "ok"; rcs.Status != status {
		t.Errorf("Status: %q != %q", rcs.Status, status)
	}
}

// Not return ok for convert with new ID
func TestConvertWithNewID(t *testing.T) {
	session, err := NewSession(Options{
		ClientID: []byte("unknown_id"),
	})

	if err != nil {
		t.Error(err)
	}

	rc, err := session.Convert("show-bieber")
	if err != nil {
		t.Error(err)
	}

	if status := "failed"; rc.Status != status {
		t.Errorf("Status: %q != %q", rc.Status, status)
	}
}

// Not return ok for convert with new experiment
func TestConvertWithNewExperiment(t *testing.T) {
	session, err := NewSession(Options{})
	if err != nil {
		t.Error(err)
	}

	rc, err := session.Convert("show-bieber")
	if err != nil {
		t.Error(err)
	}

	if status := "failed"; rc.Status != status {
		t.Errorf("Status: %q != %q", rc.Status, status)
	}
}

// Not allow bad experiment names
func TestNotAllowBadExperimentNames(t *testing.T) {
	session, err := NewSession(Options{})
	if err != nil {
		t.Error(err)
	}

	_, err = session.Participate("%%", []string{"trolled", "not-trolled"}, "")
	if err == nil {
		t.Error("%% is not allowed as experiment names")
	}
}

// Not allow single alternative name
func TestNotAllowBadSingleAlternativeName(t *testing.T) {
	session, err := NewSession(Options{})
	if err != nil {
		t.Error(err)
	}

	_, err = session.Participate("show-bieber", []string{"trolled"}, "")
	if err == nil {
		t.Error("Alternative names must at least 2 items")
	}
}

// Not allow bad alternative names
func TestNotAllowBadAlternativeNames(t *testing.T) {
	session, err := NewSession(Options{})
	if err != nil {
		t.Error(err)
	}

	_, err = session.Participate("show-bieber", []string{"trolled", "%%"}, "")
	if err == nil {
		t.Error("%% is not allowed as alternative names")
	}
}

// It works
func TestWorkflow(t *testing.T) {
	session, err := NewSession(Options{})
	if err != nil {
		t.Error(err)
	}

	rc, err := session.Convert("testing")
	if err != nil {
		t.Error(err)
	}

	if status := "failed"; rc.Status != status {
		t.Errorf("Status: %q != %q", rc.Status, status)
	}

	pOne, err := session.Participate("testing", []string{"one", "two"}, "")
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 2; i++ {
		rl, err := session.Participate("testing", []string{"one", "two"}, "")
		if err != nil {
			t.Error(err)
		}

		if rl.Alternative.Name != pOne.Alternative.Name {
			t.Errorf("Alternative name: %q != %q", rl.Alternative.Name, pOne.Alternative.Name)
		}
	}

	rOne, err := session.Convert("testing")
	if err != nil {
		t.Error(err)
	}

	if status := "ok"; rOne.Status != status {
		t.Errorf("Status: %q != %q", rOne.Status, status)
	}

	oldClientID := session.Options.ClientID

	id, err := GenerateClientID()
	if err != nil {
		t.Error(err)
	}

	session.Options.ClientID = id
	nOne, err := session.Convert("testing")
	if err != nil {
		t.Error(err)
	}

	if status := "failed"; nOne.Status != status {
		t.Errorf("Status: %q != %q", nOne.Status, status)
	}

	pTwo, err := session.Participate("testing", []string{"one", "two"}, "")
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 2; i++ {
		rl, err := session.Participate("testing", []string{"one", "two"}, "")
		if err != nil {
			t.Error(err)
		}

		if rl.Alternative.Name != pTwo.Alternative.Name {
			t.Errorf("Alternative name: %q != %q", rl.Alternative.Name, pTwo.Alternative.Name)
		}
	}

	nTwo, err := session.Convert("testing")
	if err != nil {
		t.Error(err)
	}

	if status := "ok"; nTwo.Status != status {
		t.Errorf("Status: %q != %q", nTwo.Status, status)
	}

	session.Options.ClientID = oldClientID
	for i := 0; i < 2; i++ {
		rl, err := session.Participate("testing", []string{"one", "two"}, "")
		if err != nil {
			t.Error(err)
		}

		if rl.Alternative.Name != pOne.Alternative.Name {
			t.Errorf("Alternative name: %q != %q", rl.Alternative.Name, pOne.Alternative.Name)
		}
	}
}
