package constants

type Constants struct {
	NotEnoughSpace string
}

func GetConstants() *Constants{
	return &Constants{
		NotEnoughSpace: "Not enough space to render pannels"
	}
}
