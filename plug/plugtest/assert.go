package plugtest

func (t *PT) AssertSuccess() {
	t.T.Helper()
	t.after()
	if t.Err != nil {
		t.T.Log(t.logbuf.String())
		t.T.Fatal("should have succeeded", t.Err)
	}
}

func (t *PT) AssertFail() {
	t.T.Helper()
	t.after()
	if t.Err == nil {
		t.T.Log(t.logbuf.String())
		t.T.Fatal("should have failed")
	}

}


func (t *PT) AssertOutput(text string) {
	t.T.Helper()
	t.after()
	// todo: document that buffer is emtpied
	out := t.Output()
	if out != text {
		t.T.Fatalf("output not as expected!\n got:\n%s\n expected:\n%s", out, text)
	}
}
