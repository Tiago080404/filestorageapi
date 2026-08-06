<script setup lang="ts">
import { onMounted, ref } from "vue";
import Datadisplayer from "./Datadisplayer.vue";
interface DirData {
  name: string;
  url: string;
  dir: boolean;
  path: string;
}
let dirs = ref<DirData[]>([]);
const selectedFile = ref("");
const fileType = ref("");

onMounted(async () => {
  await getData();
});
const getData = async () => {
  const response = await fetch("http://localhost:8080/api/list", {
    method: "GET",
    //credentials: "include",
  });
  dirs.value = await response.json();
  console.log(dirs);
};
const getDir = async (path: string) => {
  const response = await fetch(`http://localhost:8080/api/list/${path}`, {
    method: "GET",
  });
  dirs.value = await response.json();
  console.log(dirs.value);
};

const selectFile = (path: string) => {
  console.log(path);
  selectedFile.value = path;
  if (selectedFile.value.includes("pdf")) {
    fileType.value = "pdf";
  } else if (
    selectedFile.value.includes("MOV") ||
    selectedFile.value.includes("MP4")
  ) {
    fileType.value = "vid";
  } else {
    fileType.value = "img";
  }
};
</script>
<template>
  <div class="flex flex-row flex-wrap gap-4">
    <div v-for="data in dirs" class="flex flex-col items-center">
      <img
        @click="getDir(data.path)"
        class="w-20 h-30"
        v-if="data.dir === true"
        src="../assets/folder.svg"
        alt=""
      />
      <img
        @click="selectFile(data.path)"
        class="w-20 h-30"
        v-else
        src="../assets/file.svg"
        alt=""
      />
      <p class="w-20 text-center text-sm truncate">{{ data.name }}</p>
    </div>

    <div v-if="selectedFile">
      <Datadisplayer
        v-if="selectedFile"
        :selected-file-name="selectedFile"
        :file-type="fileType"
        @close="selectFile('')"
      ></Datadisplayer>
    </div>
  </div>
</template>
