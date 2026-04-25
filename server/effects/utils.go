package effects



type EffectFunc func(inputImage []byte) (outputImage []byte, err error)