<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import Datadisplayer from "./Datadisplayer.vue";
import FileUploader from "./FileUploader.vue";
import Navbar from "./Navbar.vue";

interface DirData {
  name: string;
  url: string;
  dir: boolean;
  path: string;
}

const selectedFile = ref("");
const fileType = ref("");
const onlyFilesData = ref<string[]>([]);
const apiUrl = import.meta.env.VITE_API_URL;

let dirs = ref<DirData[]>([]);
let fileIdx = ref(0);
let fileAdd = ref(false);
let currentPath = ref("");

onMounted(async () => {
  await getData();
  getOnlyFiles();
});

const getData = async () => {
  currentPath.value = "";
  const response = await fetch(`${import.meta.env.VITE_API_URL}api/list`, {
    method: "GET",
    //credentials: "include",
  });
  dirs.value = await response.json();
  console.log(dirs.value);
};

const getDir = async (path: string) => {
  console.log("getDir:", path);
  currentPath.value = path;
  const response = await fetch(
    `${import.meta.env.VITE_API_URL}api/list/${path}`,
    {
      method: "GET",
    },
  );
  dirs.value = await response.json();
  console.log(dirs.value);

  fileIdx.value = 0;
  getOnlyFiles();
};

const selectFile = (path: string) => {
  selectedFile.value = path;
  fileIdx.value = onlyFilesData.value.indexOf(path);
  console.log(selectedFile.value)
  if (selectedFile.value.includes("pdf")) {
    fileType.value = "pdf";
  } else if (
    selectedFile.value.includes("mov") ||
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
const closeViewer = async () => {
  selectedFile.value = "";
  fileType.value = "";
  fileIdx.value = 0;
  if (dirs.value[0].path.includes("/")) {
    await getDir(currentPath.value);
  } else {
    await getData();
  }
};
const closeUploader = async () => {
  fileAdd.value = false;
  if (dirs.value[0].path.includes("/")) {
    await getDir(currentPath.value);
  } else {
    await getData();
  }
};
const goFolder = async () => {
  let pos = currentPath.value.lastIndexOf("/");
  if (pos === -1) {
    //homedir
    await getData();
  } else {
    console.log(currentPath.value.slice(0, pos));
    currentPath.value = currentPath.value.slice(0, pos);
    await getDir(currentPath.value);
  }
};
</script>
<template>
  <div>
    <div class="fixed bottom-0 left-0 right-0 z-50  rounded-lg bg-[#C5AB9F] shadow-md">
      <div class="flex justify-center">

        <img
          @click="fileAdd = true"
          src="../assets/uploadfile.svg"
          class="w-15 h-15"
        />
        <Navbar v-if="currentPath !== ''" @back="goFolder"></Navbar>
      </div>
    </div>
  </div>

  <div class="flex flex-row flex-wrap gap-4">
    <div class="flex flex-col items-center justify-center"></div>
    <div v-for="data in dirs" class="flex flex-col items-center justify-center">
      <img
        @click="getDir(data.path)"
        class="w-20 h-30"
        v-if="data.dir === true"
        src="../assets/folder.svg"
        alt=""
      />
      <img
        @click="selectFile(data.path)"
        class="w-20 h-20 p-2"
        v-else
        :src="`${apiUrl}api/preview/${data.name}`"
        loading="lazy"
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
      <FileUploader @close="closeUploader"></FileUploader>
    </div>
  </div>
</template>
<!-- vllt ein button auf der aktullen dir und dann kann man da selecten was man downloaden will oder ganzer ordner downloaden -->
