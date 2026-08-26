<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import Datadisplayer from "./Datadisplayer.vue";
import FileUploader from "./FileUploader.vue";

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
let fileAdd = ref(false);

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
  fileIdx.value = onlyFilesData.value.indexOf(path);

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
    <img
      @click="fileAdd = true"
      src="../assets/uploadfile.svg"
      class="w-10 h-4 sm:w-4 sm:h-10 md:w-10 md:h-9"
    />
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
      <div
        class="fixed z-50 right-1 sm:right-2 md:right-4 bottom-4 sm:bottom-10 md:bottom-20 flex flex-col items-center"
      >
        <img
          v-if="!checkLastItem"
          @click="next"
          src="../assets/next-svgrepo-com.svg"
          class="w-10 h-14 sm:w-14 sm:h-20 md:w-20 md:h-30"
          alt=""
        />
      </div>

      <div
        class="fixed z-50 left-1 sm:left-2 md:left-4 bottom-4 sm:bottom-10 md:bottom-20 flex flex-col items-center"
      >
        <img
          v-if="!checkFirstItem"
          @click="backwards"
          src="../assets/back-light-svgrepo-com.svg"
          class="w-10 h-14 sm:w-14 sm:h-20 md:w-20 md:h-30"
          alt=""
        />
      </div>
    </div>
    <div v-if="fileAdd">
      <FileUploader></FileUploader>
    </div>
  </div>
</template>
