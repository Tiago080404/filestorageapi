<script setup lang="ts">
const emit = defineEmits(["back"]);

const moveBack = () => {
  emit("back");
};
</script>
<template>
      <img
        src="../assets/back-light-svgrepo-com.svg"
        alt=""
        class="w-15 h-15 "
        @click="moveBack"
      />

</template>
