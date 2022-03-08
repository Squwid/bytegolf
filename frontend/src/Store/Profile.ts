import { BackendURL } from "../Globals";
import { BasicProfile, Claims } from "../Types";

export const GetProfile = async(bgid: string): Promise<BasicProfile|undefined> => {
  return new Promise((res, rej) => {
    let url = `${BackendURL()}/profile/${bgid}`;

    console.log(`*** ${url}`);

    fetch(url, {credentials: 'include'})
      .then(resp => {
        if (resp.status === 404) return res(undefined);
        if (resp.status !== 200) return rej(`Got bad status code ${resp.status}`);
        res(resp.json());
      })
      .catch(err => rej(`Error getting profile : ${err}`));
  })
}

export const GetClaims = async(): Promise<Claims|undefined> => {
  return new Promise((res, rej) => {
    const url = `${BackendURL()}/claims`;
    console.log(`*** ${url}`);

    fetch(url, {credentials: 'include'})
      .then(resp => {
        if (resp.status === 401) return res(undefined);
        if (resp.status !== 200) return rej(`Got bad status code ${resp.status}`);
        res(resp.json());
      })
      .catch(err => rej(`Error getting claims ${err}`));
  })
}