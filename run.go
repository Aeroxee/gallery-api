package galleryapi

func Run(addr string) {
	r := router()
	r.Run(addr)
}
