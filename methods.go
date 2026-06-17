package ksi

func (k *ksi) Get(pattern string, handler KsiFunc) {
	k.Handle("GET "+pattern, handler)
}

func (k *ksi) Post(pattern string, handler KsiFunc) {
	k.Handle("POST "+pattern, handler)
}

func (k *ksi) Put(pattern string, handler KsiFunc) {
	k.Handle("PUT "+pattern, handler)
}

func (k *ksi) Patch(pattern string, handler KsiFunc) {
	k.Handle("PATCH "+pattern, handler)
}

func (k *ksi) Delete(pattern string, handler KsiFunc) {
	k.Handle("DELETE "+pattern, handler)
}

func (k *ksi) Head(pattern string, handler KsiFunc) {
	k.Handle("HEAD "+pattern, handler)
}

func (k *ksi) Options(pattern string, handler KsiFunc) {
	k.Handle("OPTIONS "+pattern, handler)
}

func (k *ksi) Connect(pattern string, handler KsiFunc) {
	k.Handle("CONNECT "+pattern, handler)
}

func (k *ksi) Trace(pattern string, handler KsiFunc) {
	k.Handle("TRACE "+pattern, handler)
}
