<template>
  <div class="roll-wrap" :style="{ fontSize: `${cellHeight}px` }">
    <ul class="roll-box">
      <li
        v-if="type === 'number'"
        class="roll-item"
        v-for="(item, index) in numberArr"
        :key="index"
        :style="{ height: `${cellHeight}px`, lineHeight: `${cellHeight}px` }"
      >
        <!--小数点或其他情况-->
        <div v-if="isNaN(parseFloat(item))">{{ item }}</div>
        <div v-else :style="getStyles(index)">
          <!--数字0到9-->
          <div
            :style="{
              height: `${cellHeight}px`,
              lineHeight: `${cellHeight}px`,
            }"
            v-for="(subItem, subIndex) in oneToNineArr"
            :key="subIndex"
          >
            {{ subItem }}
          </div>
        </div>
      </li>
      <li
        v-if="type === 'letter'"
        class="roll-item"
        v-for="(item, index) in wordLetters"
        :key="index"
        :style="{ height: `${cellHeight}px`, lineHeight: `${cellHeight}px` }"
      >
        <div :style="getStyles(index)">
          <!--数字0到9-->
          <div
            :style="{
              display: 'block',
              letterSpacing: '-1px',
              height: `${cellHeight}px`,
              lineHeight: `${cellHeight}px`,
            }"
            v-for="(subItem, subIndex) in AtoZLetters"
            :key="subIndex"
          >
            {{ subItem }}
          </div>
        </div>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  name: "TextRoll",
  props: {
    type: String,
    // Height, default 30
    cellHeight: {
      type: Number,
      default: 30,
    },
    // Rolling number that needs to be passed in
    rollNumber: {
      type: [String, Number],
      default: 0,
    },
    rollWord: {
      type: String,
    },
    // Scroll duration, unit ms. Default 1.5s
    dur: {
      type: Number,
      default: 1500,
    },
    // Easing function, default ease
    easeFn: {
      type: String,
      default: "ease",
    },
  },
  data() {
    const { rollNumber, rollWord } = this;
    return {
      // Incoming number
      number: `${rollNumber}`,
      // The number passed in is parsed as an array
      numberArr: [],
      // Offset
      numberOffsetArr: [],
      // 0 to 9 array
      oneToNineArr: [0, 1, 2, 3, 4, 5, 6, 7, 8, 9],

      //Incoming word
      word: `${rollWord}`,
      wordLetters: [],
      wordLettersOffsetArr: [],
      // AtoZLetters: ['A','B','C','D','E','F','G','H','I','J','K','L','M','N','O','P','Q','R','S','T','U','V','W','X','Y','Z']
      AtoZLetters: ["A", "B", "C", "E", "H", "I", "N", "T", "X"],
    };
  },
  created() {
    console.log(this.$props.type);
    if (this.$props.type === "number") {
      this.numberArr = this.number.split("");
      this.resetState(this.numberArr.length);
    }
    if (this.$props.type === "letter") {
      this.wordLetters = this.word.split("");
      this.resetState(this.wordLetters.length);
    }
  },
  watch: {
    rollNumber(value, oldVal) {
      this.number = `${value}`;
      this.numberArr = `${value}`.split("");
      this.resetState(this.numberArr.length);
    },
    rollWord(value, oldVal) {
      this.word = `${value}`;
      this.wordLetters = `${value}`.split("");
      this.resetState(this.wordLetters.length);
    },
  },
  methods: {
    resetState(len) {
      const newArr = new Array(len).join(",").split(",");
      switch (this.$props.type) {
        case "number":
          this.numberOffsetArr = newArr.map(() => 0);
          // Delayed execution of animation
          setTimeout(() => {
            // Set the corresponding offset of the incoming digital subscript,Reassign
            this.numberArr.forEach((num, i) => {
              this.numberOffsetArr[i] = num * this.cellHeight;
            });
          }, 30);
          break;
        case "letter":
          this.wordLettersOffsetArr = newArr.map(() => 0);
          // Delayed execution of animation
          setTimeout(() => {
            // Set the corresponding offset of the incoming digital subscript and re-assign it
            this.wordLetters.forEach((letter, i) => {
              let offUnit = this.getOffsetUnit(letter);
              this.wordLettersOffsetArr[i] = offUnit * this.cellHeight;
            });
          }, 30);
          break;
      }
    },

    getStyles(index) {
      return this.$props.type === "number"
        ? {
            transition: `${this.easeFn} ${this.dur}ms`,
            transform: `translate(0%, -${this.numberOffsetArr[index]}px)`,
          }
        : {
            padding: 0,
            transition: `${this.easeFn} ${this.dur}ms`,
            transform: `translate(0%, -${this.wordLettersOffsetArr[index]}px)`,
          };
    },
    getOffsetUnit(letter) {
      return this.AtoZLetters.indexOf(letter);
    },
  },
};
</script>
<style lang="scss" scoped>
.roll-wrap {
  ul.roll-box {
    display: flex;
    padding: 0;
    margin: 0;
    text-align: center;
    overflow: hidden;
    li {
      overflow: hidden;
    }
  }
}
</style>
