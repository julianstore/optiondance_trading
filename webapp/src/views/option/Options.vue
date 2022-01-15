<template>
  <div class="trade-info">
    <back-header back-to="/home"></back-header>
    <div class="tab" v-bind:class="{ 'end-tab': status === TradeState.END }">
      <ul class="status-filter">
        <li
          @click="changeStatus(TradeState.MATCHING)"
          v-bind:class="{ selected: status === TradeState.MATCHING }"
        >
          {{ $t("options.statusMatching") }}
        </li>
        <li
          @click="changeStatus(TradeState.POSITION)"
          v-bind:class="{ selected: status === TradeState.POSITION }"
        >
          {{ $t("options.statusOpen") }}
        </li>
        <li
          @click="changeStatus(TradeState.END)"
          v-bind:class="{ selected: status === TradeState.END }"
        >
          {{ $t("options.statusEnd") }}
        </li>
      </ul>
      <ul class="type-filter" v-show="status === 'end'">
        <li
          @click="changeFilterType('order')"
          v-bind:class="{ selected: type === 'order' }"
        >
          {{ $t("options.order") }}
        </li>
        <li
          @click="changeFilterType('position')"
          v-bind:class="{ selected: type === 'position' }"
        >
          {{ $t("options.contract") }}
        </li>
      </ul>
    </div>

    <div
      ref="wrapper"
      class="wrapper"
      v-bind:style="{
        marginTop: (status === TradeState.END ? 154 : 114) + 'px',
      }"
    >
      <div
        class="content"
        v-bind:style="{ height: contentHeightString + 'px' }"
      >
        <!--订单-->
        <div v-if="status === TradeState.MATCHING" ref="content">
          <div
            v-if="openOrders !== undefined && openOrders.length > 0"
            v-for="(e, index) in openOrders"
            class="order-item"
          >
            <OrderCard
              :id="e.order_id"
              :order-num="e.order_num_string"
              :remaining-amount="e.remaining_amount"
              :remaining-funds="e.remaining_funds"
              :filled-amount="e.filled_amount"
              :filled-funds="e.filled_funds"
              :expiration-date="e.expiration_date"
              :option-type="e.option_type"
              :side="e.side"
              :instrument-name="e.instrument_name"
              :order-status="e.order_status"
              :price="e.price"
              :strike-price="e.strike_price"
              :base-currency="e.base_currency"
              :quote-currency="e.quote_currency"
              :margin="e.margin"
              :clickable="true"
            />
          </div>
          <div v-else class="empty-tip">
            <span>{{ $t("options.noOrder") }}</span>
          </div>
        </div>
        <!--positions-->
        <div v-if="status === TradeState.POSITION" ref="content">
          <div
            v-if="openPositions !== undefined && openPositions.length > 0"
            v-for="(e, index) in openPositions"
            class="order-item"
          >
            <PositionCard
              :clickable="true"
              :size="e.size"
              :id="e.position_id"
              :num="e.position_num_string"
              :expiration-date="e.expiration_date"
              :option-type="e.type"
              :side="e.side"
              :funds="e.funds"
              :instrument-name="e.instrument_name"
              :exercised-size="e.exercised_size"
              :status="e.status"
              :price="e.price"
              :display-id="false"
              :instrument="e.instrument"
              :quote-currency="e.quote_currency"
              :base-currency="e.base_currency"
              :position-funds="e.initial_funds"
            />
          </div>
          <div v-else class="empty-tip">
            <span>{{ $t("options.noContract") }}</span>
          </div>
        </div>
        <!--已结束-->
        <div v-if="status === TradeState.END" ref="content">
          <div v-if="type === 'order'">
            <div
              v-if="closedOrders !== undefined && closedOrders.length > 0"
              v-for="(e, index) in closedOrders"
              class="order-item"
            >
              <OrderCard
                :id="e.order_id"
                :order-num="e.order_num_string"
                :display-id="false"
                :remaining-amount="e.remaining_amount"
                :remaining-funds="e.remaining_funds"
                :filled-amount="e.filled_amount"
                :filled-funds="e.filled_funds"
                :expiration-date="e.expiration_date"
                :option-type="e.option_type"
                :side="e.side"
                :instrument-name="e.instrument_name"
                :order-status="e.order_status"
                :price="e.price"
                :strike-price="e.strike_price"
                :base-currency="e.base_currency"
                :quote-currency="e.quote_currency"
                :margin="e.margin"
                :clickable="true"
              />
            </div>
            <div v-else class="empty-tip">
              <span>{{ $t("options.noOrder") }}</span>
            </div>
          </div>
          <div v-if="type === 'position'">
            <div
              v-if="closedPositions !== undefined && closedPositions.length > 0"
              v-for="(e, index) in closedPositions"
              class="order-item"
            >
              <PositionCard
                :clickable="true"
                :size="e.size"
                :id="e.position_id"
                :num="e.position_num_string"
                :expiration-date="e.expiration_date"
                :option-type="e.type"
                :side="e.side"
                :funds="e.funds"
                :instrument-name="e.instrument_name"
                :status="e.status"
                :exercised-size="e.exercised_size"
                :price="e.price"
                :display-id="false"
                :instrument="e.instrument"
                :quote-currency="e.quote_currency"
                :base-currency="e.base_currency"
                :position-funds="e.initial_funds"
              />
            </div>
            <div v-else class="empty-tip">
              <span>{{ $t("options.noContract") }}</span>
            </div>
          </div>
        </div>
      </div>
      <div class="scroll-loading" v-show="scrollLoadingVisible">
        <div id="scroll-loader"></div>
      </div>
    </div>
  </div>
</template>

<script>
import BackHeader from "@/components/layout/header/BackHeader.vue";
import { TradeStatus } from "@/util/constants";
import OrderCard from "@/components/options/OrderCard.vue";
import { setHtmlMeta } from "@/util/utils";
import PositionCard from "@/components/options/PositionCard.vue";
import $axios from "@/api";
import { mapGetters } from "vuex";
import { parseInstrumentName } from "@/util/instrument";
import BScroll from "better-scroll";
import lottie from "lottie-web";
import TabLoadingJson from "@/assets/lottie/tab_loading.json";

export default {
  name: "options",
  components: { PositionCard, OrderCard, BackHeader },
  async mounted() {
    document.body.style.backgroundColor = "#F7F7F7";
    setHtmlMeta("theme-color", "#F7F7F7");
    this.TradeState = TradeStatus;
    let status = this.$route.query.status;
    let type = this.$route.query.type;
    if (status) {
      this.status = status;
    } else {
      this.status = TradeStatus.MATCHING;
    }
    if (type) {
      this.type = type;
    }
    await this.getPageData("tab");
    await this.$nextTick(() => {
      this.initScroll();
    });
    lottie.loadAnimation({
      container: document.getElementById("scroll-loader"),
      renderer: "svg",
      loop: true,
      autoplay: true,
      animationData: TabLoadingJson,
    });
  },
  computed: {
    ...mapGetters({
      mixinToken: "user/mixinToken",
      mixinUser: "user/mixinUser",
    }),
    contentHeightString() {
      if (this.contentHeight) {
        return this.contentHeight;
      } else {
        return null;
      }
    },
  },
  data() {
    return {
      TradeState: {},
      status: "",
      type: "",
      openOrders: [],
      openPositions: [],
      closedOrders: [],
      closedPositions: [],
      scrollState: true, //Can slide
      indexScrollTop: 0,
      contentHeight: 0,
      scrollLoadingVisible: false,
    };
  },
  methods: {
    async changeStatus(status) {
      this.status = status;
      let query = {
        status: status,
      };
      if (status === TradeStatus.END) {
        let type = this.type ? this.type : "order";
        this.type = type;
        query = Object.assign(query, {
          type: type,
        });
      }
      await this.getPageData("tab");
      await this.$router.push({ name: "options", query: query });
      await this.$nextTick(() => {
        this.scroll.refresh();
      });
    },
    async changeFilterType(type) {
      this.type = type;
      await this.getPageData("tab");
      await this.$router.push({
        name: "options",
        query: {
          status: this.status,
          type: type,
        },
      });
      await this.$nextTick(() => {
        this.scroll.refresh();
        console.log("changeFilterType init");
      });
    },
    async getPageData(animType) {
      this.startAnim(animType);
      switch (this.status) {
        case TradeStatus.MATCHING:
          let res = await $axios.get("/v1/orders", {
            params: { status: "open", current: 1, size: 50 },
          });
          this.openOrders = res.data.records;
          break;
        case TradeStatus.POSITION:
          let openPositionsResult = await $axios.get("/v1/positions", {
            params: { status: "open", current: 1, size: 50 },
          });
          this.openPositions = openPositionsResult.data.records;
          for (let i = 0; i < this.openPositions.length; i++) {
            this.openPositions[i].instrument = parseInstrumentName(
              this.openPositions[i].instrument_name
            );
          }
          break;
        case TradeStatus.END:
          let type = this.type ? this.type : "order";
          this.type = type;
          if (this.type === "order") {
            let res = await $axios.get("/v1/orders", {
              params: { status: "closed", current: 1, size: 100 },
            });
            this.closedOrders = res.data.records;
            break;
          }
          if (this.type === "position") {
            let closedPositionsResult = await $axios.get("/v1/positions", {
              params: {
                status: "closed",
                current: 1,
                size: 100,
                t: new Date().getTime(),
              },
            });
            this.closedPositions = closedPositionsResult.data.records;
            for (let i = 0; i < this.closedPositions.length; i++) {
              this.closedPositions[i].instrument = parseInstrumentName(
                this.closedPositions[i].instrument_name
              );
            }
          }
          break;
      }
      this.endAnim(animType);
      await this.$nextTick(() => {
        this.contentHeight = this.$refs.content.offsetHeight + 32;
      });
    },
    startAnim(animType) {
      if (animType === "tab") {
        this.$tabLoading.show();
      }
      if (animType === "scroll") {
        this.scrollLoadingVisible = true;
      }
    },
    endAnim(animType) {
      if (animType === "tab") {
        this.$tabLoading.hide();
      }
      if (animType === "scroll") {
        this.scrollLoadingVisible = false;
      }
    },
    async initScroll() {
      this.scroll = new BScroll(this.$refs.wrapper, {
        // Pull up loading
        pullUpLoad: {
          // The pullingUp event is triggered when the pull-up distance exceeds 30px
          threshold: -30,
        },
        // Pull down to refresh
        pullDownRefresh: {
          // Pulling down more than 30px triggers the pullingDown event
          threshold: 30,
          // The rebound stays 20px from the top
          stop: 40,
        },
        mouseWheel: true,
        click: true,
      });
      this.scroll.on("pullingDown", async () => {
        this.getPageData("scroll").then(() => {
          this.$nextTick(() => {
            this.scroll.refresh(); // After the DOM structure changes, reinitialize BScroll
          });
          setTimeout(() => {
            // When things are done, you need to call this method to tell better-scroll that the data has been loaded,
            // otherwise the drop-down event will only be executed once
            this.scroll.finishPullDown();
            // setTimeout(()=>this.scroll = this.initScroll(),0)
          }, 100);
        });
      });
    },
  },
};
</script>

<style scoped lang="scss">
.empty-tip {
  margin-top: 216px;
  text-align: center;
  span {
    //width: 108px;
    height: 29px;
    background: rgba(0, 0, 0, 0.02);
    border-radius: 27px;
    padding: 8px 26px;
    font-size: 14px;
    line-height: 21px;
    text-align: center;
    color: #88888b;
    box-sizing: border-box;
    margin: 0 auto;
  }
}
.trade-info {
  .tab {
    position: fixed;
    top: 54px;
    background: #f7f7f7;
    z-index: 1200;
    padding: 13px 25px 16px;
    margin-bottom: 10px;
    box-sizing: border-box;
    width: 100%;
    ul.status-filter {
      display: flex;
      flex-direction: row;
      li {
        margin-right: 36px;
        font-style: normal;
        font-weight: bold;
        font-size: 20px;
        line-height: 30px;
        color: #b4b4b6;
        &.selected {
          color: #151516;
          position: relative;
          //&:after {
          //  content: ' ';
          //  position: absolute;
          //  right: -8px;
          //  top: -2px;
          //  width: 5px;
          //  height: 5px;
          //  background-color: #FF6F61;;
          //  border-radius: 2.5px;
          //}
        }
      }
    }
    .type-filter {
      margin-top: 14px;
      display: flex;
      flex-direction: row;
      li {
        margin-right: 10px;
        display: flex;
        align-items: center;
        justify-content: center;
        width: 72px;
        height: 28px;
        background: rgba(136, 136, 139, 0.1);
        border-radius: 4px;
        font-style: normal;
        font-weight: normal;
        font-size: 13px;
        line-height: 19px;
        text-align: center;
        color: #000000;
      }
      li.selected {
        background: #151516;
        color: #ffffff;
      }
    }
  }
  .end-tab {
    padding-bottom: 12px;
  }
  .content {
    padding: 0 24px;
    .order-item {
      margin-bottom: 20px;
    }
  }
}
.wrapper {
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  right: 0;
  overflow: hidden;
  .content {
    padding-top: 10px;
  }
  .scroll-loading {
    position: absolute;
    width: 100%;
    top: 0;
    display: flex;
    justify-content: center;
    align-items: center;
    transition: all;
    #scroll-loader {
      margin-top: 10px;
      width: 20px;
      height: 20px;
    }
  }
}
.wrapper-content {
}
</style>
