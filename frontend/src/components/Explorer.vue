<script setup lang="ts">
import { onMounted, ref } from "vue";
interface DirData {
  name: string;
  url: string;
  dir: boolean;
  path: string;
}
let dirs = ref<DirData[]>([]);
const selectedFile = ref("");
const selectedVideo = ref("");
const selectedImage = ref("");
const selectedImageName = ref("");

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
const openFile = async (file: string) => {
  if (file.includes("pdf")) {
    console.log("pdf");
    selectedFile.value = `http://localhost:8080/api/open/${file}`;
    selectedImage.value = "";
    selectedVideo.value = "";
  } else if (file.includes("MOV")) {
    selectedVideo.value = `http://localhost:8080/api/open/${file}`;
    selectedImage.value = "";
    selectedFile.value = "";
  } else {
    selectedImage.value = `http://localhost:8080/api/open/${file}`;
    selectedImageName.value = file;
    selectedVideo.value = "";
    selectedFile.value = "";
  }
};
const downloadFile = async (file: string) => {
  const response = await fetch(`http://localhost:8080/api/download/${file}`, {
    method: "GET",
  });
  const data = await response.blob();
  const url = window.URL.createObjectURL(data);
  const a = document.createElement("a");
  a.href = url;
  a.setAttribute("download",file)
  document.body.appendChild(a);
  a.click();
   URL.revokeObjectURL(url);
  document.body.removeChild(a); 
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
        @click="openFile(data.path)"
        class="w-20 h-30"
        v-else
        src="../assets/file.svg"
        alt=""
      />
      <p class="w-20 text-center text-sm truncate">{{ data.name }}</p>
    </div>
    <div v-if="selectedImage">
     
       <button @click="downloadFile(selectedImageName)">Download</button> 
      <img v-if="selectedImage" :src="selectedImage" width="600" height="400" />
    </div>
    <video
      v-if="selectedVideo"
      :src="selectedVideo"
      controls
      height="250"
    ></video>
    <iframe
      v-if="selectedFile"
      :src="selectedFile"
      style="border: none"
      width="100%"
      height="600px"
    ></iframe>
  </div>
</template>
