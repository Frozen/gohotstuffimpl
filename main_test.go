package main

type ViewNumber = int

type PublicKey = string

type QuorumCertificate struct {
	Type       string
	viewNumber ViewNumber
	Node       PublicKey
}

//type QC = QuorumCertificate

type msg struct {
	Type       string
	Node       PublicKey
	Justify    QuorumCertificate
	viewNumber ViewNumber
	partialSig string
}

type leaf struct {
	parent string
	cmd    string
}

// --------------------------------------- Functions ---------------------------------------
func Msg(t string, n string, qc QuorumCertificate) msg {
	return msg{
		Type:    t,
		Node:    n,
		Justify: qc,
	}
}

func voteMsg(t string, n string, qc QuorumCertificate) msg {
	m := Msg(t, n, qc)
	m.partialSig = tsignr(m.Type, m.viewNumber, m.node)
	return m
}

func createLeaf(parent string, cmd string) leaf {
	b := leaf{}
	b.parent = parent
	b.cmd = cmd
	return b
}

func QC(v V) QuorumCertificate {
	/*
		qc.Type  ← m.type : m ∈ V
		qc.viewNumber ← m.viewNumber : m ∈ V
		qc.node ← m.node : m ∈ V
		 qc.sig ← tcombine(hqc.type, qc.viewNumber , qc.nodei, {m.partialSig | m ∈ V })
		 return qc
	*/
	QuorumCertificate{
		Type: v.Type, // {$new-view, prepare, pre-commit, commit, decide$}
		viewNumber: v.ViewNumber

	}
}

func matchingMsg(m msg, t string, v ViewNumber) bool {
	return (m.Type == t) && (m.viewNumber == v)
}

func matchingQC(qc QC, t string, v ViewNumber) bool {
	return (qc.Type == t) && (qc.viewNumber == v)
}
