import TabLoading from "./TabLoading.vue";
import {createComponent,removeElement} from "@/components/common/utils";

export function useTabLoading(globalProps = {}, globalSlots = {}) {
    let instance =  null;
    return {
        show(props = globalProps, slots = globalSlots) {
            const forceProps = {
                programmatic: true,
                lockScroll: true,
                isFullPage: false
            };
            const propsData = Object.assign({}, globalProps, props, forceProps);
            let container = propsData.container;

            if (!propsData.container) {
                container = document.body;
                propsData.isFullPage = true;
            }
            instance = createComponent(TabLoading, propsData, container);
            const mergedSlots = Object.assign({}, globalSlots, slots);
            Object.keys(mergedSlots).map((name) => {
                if (instance != null) {
                    instance.slots[name] = mergedSlots[name]
                }
            });
        },
        hide() {
            if (instance != null) {
                let root = instance.vnode.el;
                removeElement(root.parentElement);
            }
        }
    }
}


const Plugin = (app, options = {}) => {
    let methods = useTabLoading()
    app.$tabLoading = methods
    app.config.globalProperties.$tabLoading = methods
}

TabLoading.install = Plugin

export default TabLoading