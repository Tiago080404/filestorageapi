<script setup lang="ts">
const emit = defineEmits(["back"]);

const moveBack = () => {
  emit("back");
};
</script>
<template>
  <div class="flex justify-center align-middle">
    <div>
      <img
        src="../assets/back-light-svgrepo-com.svg"
        alt=""
        class="w-15 m-2"
        @click="moveBack"
      />
    </div>
  </div>
</template>
