import {h, render} from "vue";

export function createComponent(component, props, parentContainer) {
    const vnode = h(component, props)
    const container = document.createElement('div')
    parentContainer.appendChild(container);
    render(vnode, container);

    return vnode.component
}

export function removeElement(el) {
    if (typeof el.remove !== 'undefined') {
        el.remove()
    } else {
        el.parentNode?.removeChild(el)
    }
}