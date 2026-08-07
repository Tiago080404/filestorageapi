<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
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
const onlyFilesData = ref<string[]>([]);
let fileIdx = ref(0);

onMounted(async () => {
  await getData();
  getOnlyFiles();
});

const getData = async () => {
  const response = await fetch("http://localhost:8080/api/list", {
    method: "GET",
    //credentials: "include",
  });
  dirs.value = await response.json();
};
const getDir = async (path: string) => {
  const response = await fetch(`http://localhost:8080/api/list/${path}`, {
    method: "GET",
  });
  dirs.value = await response.json();
  fileIdx.value = 0;
  getOnlyFiles();
};

const selectFile = (path: string) => {
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

  console.log("path in explorer", path, "idx", fileIdx.value);
};

const getOnlyFiles = () => {
  onlyFilesData.value = [];
  for (let i = 0; i < dirs.value.length; i++) {
    if (dirs.value[i].dir === false) {
      onlyFilesData.value.push(dirs.value[i].path);
    }
  }
};
const next = () => {
  if (fileIdx.value === onlyFilesData.value.length - 1) {
    return;
  }
  fileIdx.value += 1;
  const dir = onlyFilesData.value[fileIdx.value];
  console.log("next: ", dir);
  selectFile(dir);
};
const backwards = () => {
  if (fileIdx.value === 0) {
    return;
  }
  fileIdx.value--;
  const dir = onlyFilesData.value[fileIdx.value];
  selectFile(dir);
};

const checkLastItem = computed(() => {
  if (fileIdx.value === onlyFilesData.value.length - 1) {
    return true;
  }
  return false;
});
const checkFirstItem = computed(() => {
  if (fileIdx.value === 0) {
    return true;
  }
  return false;
});
const closeViewer = () => {
  selectedFile.value = "";
  fileType.value = "";
  fileIdx.value = 0;
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
        :key="selectedFile"
        v-if="selectedFile"
        :selectedFilePath="selectedFile"
        :file-type="fileType"
        @close="closeViewer"
      ></Datadisplayer>
      <div class="fixed z-50 right-2 bottom-100">
        <button @click="next" :class="{ disable: checkLastItem }">
          Forward
        </button>
      </div>
      <div class="fixed z-50 left-2 bottom-100">
        <button @click="backwards" :class="{ disable: checkFirstItem }">
          Backwards
        </button>
      </div>
    </div>
  </div>
</template>
<style>
.disable {
  color: blue;
  font-weight: bold;
}
</style>
