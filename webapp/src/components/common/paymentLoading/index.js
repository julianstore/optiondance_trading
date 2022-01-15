import PaymentLoading from "./PaymentLoading.vue";
import {createComponent,removeElement} from "@/components/common/utils";

export function useLoading(globalProps = {}, globalSlots = {}) {
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
            instance = createComponent(PaymentLoading, propsData, container);
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
    let methods = useLoading()
    app.$paymentLoading = methods
    app.config.globalProperties.$paymentLoading = methods
}

PaymentLoading.install = Plugin

export default PaymentLoading