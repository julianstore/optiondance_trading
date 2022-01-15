<template>
  <div class="chart">
  <vue3-chart-js v-if="type==='premium'" v-bind="{...premiumChartData}" />
  <vue3-chart-js v-if="type==='underlying'" v-bind="{...underlyingChartData}" />
<!--    <vue3-chart-js v-bind="{ ...underlyingChart }" />-->
  </div>
</template>

<script>
import Vue3ChartJs from "@j-t-mcc/vue3-chartjs";
import {QuoteCurrency} from "@/util/constants";

export default {
  name: "IncomeChart",
  components: {
    Vue3ChartJs,
  },
  props:{
    dataSet:Array,
    type:String,
    year: String,
  },
  setup(props) {
    const getMaxAmount = (dataSet) => {
      let max = 0
      for (let i = 0; i < dataSet.length; i++) {
        if (dataSet[i] > max)
          max = Number(dataSet[i])
      }
      if (max < 100) {
        max = 100
      }else if(max <= 500) {
        max = 500
      }else {
        max = (max / 500).toFixed(0) * 1000 .toFixed(0)
      }
    }

    const getMinAmount = (dataSet) => {
      let min = 0
      for (let i = 0; i < dataSet.length; i++) {
        if (dataSet[i] < min)
          min = Number(dataSet[i])
      }
      if (min >= 0 ) min = 0
      else if (min > -100) min = -100
      else if (max >= -500) min = -500
      else min = (min / 500).toFixed(0) * 1000 .toFixed(0)
    }


    let month = ['January','February','March','April','May','June','July','August','September','October','November','December']
    let max = getMaxAmount(props.dataSet)
    let min = getMinAmount(props.dataSet)
    const premiumChartData = {
      type: "bar",
      options: {
        min: 0,
        max: max * 1.2,
        aspectRatio:1,
        responsive: true,
        plugins: {
          legend: {
            display: false,
            position: "top",
          },
          tooltip:{
            callbacks: {
              title(x) {
                return month[x[0].label-1] + " "+ props.year
              },
              label: function(context) {
                let label = ''
                if (context.parsed.y !== null) {
                  label += context.parsed.y + QuoteCurrency
                }
                return label;
              }
            }
          }
        },
        scales: {
          y: {
            min: min,
            max: max,
            ticks: {
              callback: function (value) {
                return `${value}`;
              },
            },
          },
        },
      },
      data: {
        labels: ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"],
        datasets: [
          {
            barPercentage:1.0,
            categoryPercentage:0.8,
            label: "OptionDance 优买收益",
            backgroundColor: ["#FF6F61"],
            data: props.dataSet,
          },
        ],
      },
    };
    const underlyingChartData = {
      type: "bar",
      options: {
        min: min,
        max: 100,
        aspectRatio:1,
        responsive: true,
        plugins: {
          legend: {
            display: false,
            position: "top",
          },
          tooltip:{
            callbacks: {
              title(x) {
                return month[x[0].label-1] + " "+ props.year
              },
              label: function(context) {
                let label = ''
                if (context.parsed.y !== null) {
                  label += context.parsed.y + 'BTC'
                }
                return label;
              }
            }
          }
        },
        scales: {
          y: {
            min: 0,
            max: max,
            ticks: {
              callback: function (value) {
                return `${value}`;
              },
            },
          },
        },
      },
      data: {
        labels: ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"],
        datasets: [
          {
            barPercentage:1.0,
            categoryPercentage:0.8,
            label: "OptionDance 优买购入BTC",
            backgroundColor: ["#FF922D"],
            data: props.dataSet,
          },
        ],
      },
    };
    return {
      premiumChartData,
      underlyingChartData
    };
  },
};
</script>
<style lang="scss">
.chart{
  display: block;
  max-width: 327px;
  margin: 0 auto 16px;
  background-color: #FFFFFF;
  padding: 20px;
  border-radius: 12px;
}
</style>