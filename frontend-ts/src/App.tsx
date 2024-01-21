import * as React from 'react';
import { DataGrid, GridRowsProp, GridValidRowModel, GridColDef } from '@mui/x-data-grid';
import { useEffect, useState } from 'react';


type ITask= {
  id: number;
  userID: number;
	taskName: string;
	taskDescription: string;
	dueDate: string;
	expirationDate: string;
}
const rows: GridRowsProp = [
  { id: 1, taskName: 'First Task', taskDescription: 'Destroy All Humans', expirationDate: '1/31' },
  { id: 2, taskName: 'Second Task', taskDescription: 'Rebuild Civilization', expirationDate: '2/31' },
];



const columns: GridColDef[] = [
  { field: 'taskName', headerName: 'Task Name', width: 150 },
  { field: 'taskDescription', headerName: 'Task Description', width: 150 },
  { field: 'expirationDate', headerName: 'Expiration Date', width: 150 },

];

export default function App() {

  const [taskData, setTaskData] = useState<ITask[]>([]);
  // const taskArray = taskData?.map((task) => <li>{task.taskName}</li>);
  console.log('declaring app');
  useEffect(() => {
    console.log('effect');
    fetch("http://localhost:8080/1001/tasks", {
      method: "GET"
    })
      .then((response) => response.json())
      .then((data) => {
        setTaskData(data);
        console.log("Got Task Data")
      })
      .catch((error) => console.log(error));
      console.log('Task Data -', taskData);
  }, []);

  const taskGridData = taskData as GridRowsProp;
  console.log('grid data - ', taskGridData);



  
  return (
    <div style={{ height: 300, width: '100%' }}>
      <DataGrid rows={taskGridData} columns={columns} />
    </div>
  );
}
