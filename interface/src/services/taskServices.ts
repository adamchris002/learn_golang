import axios from "axios";
import dayjs from "dayjs";
import customParseFormat from 'dayjs/plugin/customParseFormat'

dayjs.extend(customParseFormat)

const api = axios.create({
  baseURL: import.meta.env.VITE_BASE_URL,
})

type task = {
  title: string;
  description: string;
  task_start: string;
  due_date: string;
  completed: boolean;
};

export type TaskResponse = {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;

  title: string;
  description: string;
  task_start: string;
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

export async function getAllTasks(userId: number) {
  try {
    const result = await api.get(
      `/allTasks?userId=${userId}`,
    );
    return {
      success: true,
      data: result.data,
      messageData: null,
    };
  } catch (error: any) {
    return {
      success: false,
      data: null,
      messageData: {
        status: error.response.status,
        messageTitle: error.response.data.messageTitle,
        message: error.response.data.message,
      },
    };
  }
}

export async function getAllIncompleteTasks(userId: number) {
  try {
    const result = await api.get(
      `/incompleteTasks?userId=${userId}`,
    );
    return {
      success: true,
      data: result.data,
      messageData: null,
    };
  } catch (error: any) {
    return {
      success: false,
      data: null,
      messageData: {
        status: error.response.status,
        messageTitle: error.response.data.messageTitle,
        message: error.response.data.message,
      },
    };
  }
}

export async function getOneWeekTasks(userId: number) {
  try {
    const result = await api.get(
      `/oneWeekTasks?userId=${userId}`,
    );

    return {
      success: true,
      data: result.data,
      messageData: null,
    };
  } catch (error: any) {
    return {
      success: false,
      data: null,
      messageData: {
        status: error.response.status,
        messageTitle: error.response.data.messageTitle,
        message: error.response.data.message,
      },
    };
  }
}

export async function getTodaysTasks(userId: number) {
  try {
    const result = await api.get(
      `/todaysTasks?userId=${userId}`,
    );

    return {
      success: true,
      data: result.data,
      messageData: null,
    };
  } catch (error: any) {
    return {
      success: false,
      data: null,
      messageData: {
        status: error.response.status,
        messageTitle: error.response.data.messageTitle,
        message: error.response.data.message,
      },
    };
  }
}

export async function postTodaysTask(data: task) {
  try {
    const result = await api({
      method: "POST",
      url: `/addTasks`,
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
    const result = await api({
      method: "PUT",
      url: `/updateTaskCompletion?id=${taskId}&userId=${userId}&data=${data}`,
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

export async function updateTaskValues(
  dueDate: string,
  description: string,
  subTask: { id: number; title: string; completed: boolean }[],
  taskId: number,
  userId: number,
) {
  try {
    const data = {
      dueDate: dueDate,
      description: description,
      subTask: subTask,
    };
    const result = await api({
      method: "PUT",
      url: `/updateTaskValues?id=${taskId}&userId=${userId}`,
      data: data,
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

export async function deleteTask(taskId: number, userId: number) {
  try {
    const result = await api({
      method: "DELETE",
      url: `/deleteTask?id=${taskId}&userId=${userId}`,
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

export async function deleteExistingSubtask(subTaskId: number, taskId: number) {
  try {
    const result = await api({
      method: "DELETE",
      url: `/deleteSubTask?id=${subTaskId}&taskId=${taskId}`,
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

export async function changeTaskStartDate(
  taskId: number,
  userId: number,
  taskStartDate: string,
) {
  try {
    const result = await api({
      method: "PUT",
      url: `/updateTaskStartDate?taskId=${taskId}&userId=${userId}`,
      data: {
        task_start: taskStartDate,
      },
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

// export async function changeTaskToPending(taskId: number, userId: number) {
//   try {
//     const result = await api({
//       method: "PUT",
//       url: `/changeTaskToPending?taskId=${taskId}&userId=${userId}`,
//     });
//     const messageData = {
//       status: result.status,
//       messageTitle: result.data.messageTitle,
//       message: result.data.message,
//     };
//     return messageData;
//   } catch (error: any) {
//     return {
//       status: error.response.status,
//       messageTitle: error.response.data.messageTitle,
//       message: error.response.data.message,
//     };
//   }
// }

export async function changeTaskToActive(taskId: number, userId: number, taskDueDate: string) {
  try {
 let dueDate = ""

    if (taskDueDate) {
      const parsedDueDate = dayjs(taskDueDate, "DD/MM/YYYY")
      const today = dayjs()

      if (!parsedDueDate.isBefore(today, "day")) {
        dueDate = taskDueDate
      }
    }

    const result = await api({
      method: "PUT",
      url: `/changeTaskToActive?taskId=${taskId}&userId=${userId}`,
      data: {
        due_date: dueDate,
      },
    })

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
