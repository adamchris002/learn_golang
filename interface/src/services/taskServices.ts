import axios from "axios";

type task = {
  title: string;
  description: string;
  dueDate: string;
  completed: boolean;
};

export type TaskResponse = {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;

  title: string;
  description: string;
  due_date: string;
  completed: boolean;

  userId: number;
  user: UserResponse;

  subtasks: SubTaskResponse[] | null;
};

export type UserResponse = {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;

  firstName: string;
  lastName: string;
  username: string;
  email: string;
  password: string;

  tasks: TaskResponse[] | null;
};

export type SubTaskResponse = {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;

  title: string;
  completed: boolean;

  taskId: number;
};

export async function getTodaysTasks(userId: number) {
  try {
    const result = await axios({
      method: "GET",
      url: `http://localhost:8080/todaysTasks?userId=${userId}`,
    });
    if (result.status === 200) {
      const data = result.data;
      //ga perlu pake message lah kalo sukses
      const messageData = null;
      return { data, messageData };
    }
  } catch (error: any) {
    const messageData = {
      status: error.response.status,
      messageTitle: error.response.data.messageTitle,
      message: error.response.data.message,
    };
    return { error, messageData };
  }
}

export async function postTodaysTask(data: task) {
  try {
    const result = await axios({
      method: "POST",
      url: "http://localhost:8080/addTasks",
      data: data,
    });
    return {
      status: result.status,
      messageTitle: result.data.messageTitle,
      message: result.data.message,
    };
  } catch (error: any) {
    return {
      status: error.response.status,
      messageTitle: error.response.data.messageTitle,
      message: error.response.data.message,
    };
  }
}

export async function updateTaskCompletion(
  taskId: number,
  userId: number,
  data: boolean,
) {
  try {
    const result = await axios({
      method: "PUT",
      url: `http://localhost:8080/updateTaskCompletion?id=${taskId}&userId=${userId}&data=${data}`,
    });

    const messageData = {
      status: result.status,
      messageTitle: result.data.messageTitle,
      message: result.data.message,
    };
    return messageData;
  } catch (error: any) {
    return {
      status: error.response.status,
      messageTitle: error.response.data.messageTitle,
      message: error.response.data.message,
    };
  }
}
