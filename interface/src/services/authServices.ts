import axios from "axios";

export function signIn(data: { username: string; password: string }) {
  axios({
    method: "POST",
    url: "http://localhost:8080/login",
    data: data,
  })
    .then((result) => {
      if (result.status === 200) {
      }
    })
    .catch((err) => {
      console.error(err);
    });
}

export function register(data: {
  firstName: string;
  lastName: string;
  username: string;
  password: string;
}) {
    axios({
        method: "POST",
        url: "http://localhost:8080/register",
        data: data,
    }).then((result) => {
        if (result.status === 200) {
            
        }
    }).catch((err) => {
        console.error(err)
    })
};
