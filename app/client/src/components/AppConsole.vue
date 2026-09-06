<template>
  <div class="terminal-container">
    <div class="terminal">
      <p v-for="log in logs" :key="log.CreationDate">
        [{{ formatDate(log.CreationDate) }}] {{ log.Data.Name }}, {{ log.Data.Album }}, {{ log.Data.Year }}, {{ log.Data.Rank }}
      </p>
    </div>
    <button class="update-button" @click="updateTerminal">UPDATE</button>
  </div>
</template>

<script>
export default {
  data() {
    return {
      logs: [],
    };
  },
  methods: {
    updateTerminal() {
      this.logs = []; // Vacía el array logs
      fetch('http://localhost:5000/logs')
        .then(response => response.json())
        .then(data => {
          this.logs = data;
        });
    },
    formatDate(dateString) {
      const datePart = dateString.split(' ')[0]; // Extrae solo la fecha
      const date = new Date(datePart);
      return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()}`;
    },
  },
};
</script>

<style scoped>
.terminal-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
}

.terminal {
  width: 80%;
  height: 40%;
  background-color: black;
  color: white;
  padding: 20px;
  border-radius: 10px;
  text-align: left;
}

.update-button {
  margin-top: 20px;
  padding: 10px 20px;
}
</style>