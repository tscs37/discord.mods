package main

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"sync"
)

func main() {
	fmt.Println("Starting up Discord.Mods API...")
	defer fmt.Println("Discord.Mods API finished startup.")
	registerMutator()
}

var eventCallbacks = map[string]*callbacks{}
var eventSync = sync.RWMutex{}

type callbacks map[string]*js.Object

func registerMutator() {
	observer := js.Global.Get("MutationObserver").New(mutatorCallback)
	target := js.Global.Get("document")
	obsConfig := map[string]bool{
		"childList": true,
		"subtree":   true,
	}
	observer.Call("observe", target, obsConfig)
	js.Global.Get("dmodsNS").Set("onEvent", registerCallback)
	js.Global.Get("dmodsNS").Set("dispatchEvent", eventDispatch)
}

func registerCallback(id, on string, cb *js.Object) {
	eventSync.Lock()
	defer eventSync.Unlock()
	if eventCallbacks == nil {
		eventCallbacks = map[string]*callbacks{}
	}
	if eventCallbacks[on] == nil {
		eventCallbacks[on] = &callbacks{}
	}
	(*eventCallbacks[on])[id] = cb
}

func eventDispatch(event string, obj *js.Object) {
	eventSync.RLock()
	defer eventSync.RUnlock()
	if evtCb, ok := eventCallbacks[event]; ok {
		if evtCb == nil {
			eventCallbacks[event] = &callbacks{}
		}
		for k, v := range *evtCb {
			v.Invoke(obj)
		}
	}
}

func mutatorCallback(mutations *js.Object) {
	for i := 0; i < mutations.Length(); i++ {
		mutation := mutations.Index(i)
		if mutation.Get("target") == js.Undefined {
			return
		}
		if mutation.Get("target").Get("getAttribute") == js.Undefined {
			return
		}
		if mutation.Get("target").Call("getAttribute", "class").String() != "null" {
			classList := mutation.Get("target").Get("classList")
			if classList.Call("contains", "title-wrap").Bool() {
				eventDispatch("channelSwitch", mutation)
			}
			if classList.Call("contains", "chat").Bool() {
				eventDispatch("channelSwitch", mutation)
			}
			if mutation.Get("target").Call("getAttribute", "class").Call("indexOf", "scroller messages").Int() != -1 {
				eventDispatch("newMessage", mutation)
			}
		}
	}
}
