import logo from './logo.svg';
import './App.css';
import React, { useState, useEffect } from 'react';
function App() {
  const [taskData, setTaskData] = useState(null);

  useEffect(() => {
    fetch("/1001/tasks", {
      method: "GET"
    })
      .then((response) => response.json())
      .then((data) => {
        setTaskData(data);
        console.log("Got Task Data")
        console.log(data);
      })
      .catch((error) => console.log(error));
  }, []);

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.js</code> and save to reload. To infinity and beyond!
        </p>
        <p>
          Task Data = {taskData}
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
      </header>
    </div>
  );
}

export default App;
