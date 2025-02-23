package engine

type Operation struct {
    Name string
    Args []string
}

type OperationList []*Operation
