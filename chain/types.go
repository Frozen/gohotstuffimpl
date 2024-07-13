package chain

const NumNodes = 4

type ViewNumber = int
type UniqueID int

type Chain struct {
	viewNumber   int
	Uniq         UniqueID
	network      Network
	tally        Tally
	viewNumberCh chan ViewNumber
	//
	prepareQC *QC
	lockedQC  *QC
}

type Msg struct {
	Type       MessageType
	ViewNumber int
	Node       UniqueID
	Justify    *QC
	//Signatures [NumNodes]bool
	Payload []byte
}

type QC struct {
	Type       MessageType
	Node       UniqueID
	ViewNumber int
}

//func NewMsg(messageType MessageType, viewNumber int, node UniqueID, qc *QC) Msg {
//	return Msg{
//		Type:       messageType,
//		Node:       node,
//		ViewNumber: viewNumber,
//		Justify:    qc,
//		Payload:    createProposal(),
//	}
//
//}
