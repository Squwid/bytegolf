import { BackendURL } from "../Globals"
import { BasicFullSubmission, BasicLeaderboardEntry, BasicShortSubmission } from "../Types";

// Returns -1 if does not exist, otherwise returns the score or rejects
export const GetBestHoleScore = async(hole: string): Promise<number> => {
  return new Promise((res, rej) => {
    let url = `${BackendURL()}/leaderboards?limit=1&hole=${hole}`;

    console.log(`*** ${url}`);

    fetch(url, {credentials: 'include'})
      .then(resp => {
        if (resp.status !== 200) return rej(`Got bad status code ${resp.status}`);

        (resp.json() as Promise<BasicLeaderboardEntry[]>).then(leaders => {
          if (leaders.length === 0) return res(-1);
          res(leaders[0].Length);
        })
      })
      .catch(err => rej(`Error getting best score ${err}`));
  });
}

export const GetLeaderboard = async(hole: string, limit: number, lang?: string): Promise<BasicLeaderboardEntry[]> => {
  return new Promise((res, rej) => {
    let url = `${BackendURL()}/leaderboards?limit=${limit}&hole=${hole}`;

    if (lang) url += `&lang=${lang}`;

    console.log(`*** ${url}`);

    fetch(url, {credentials: 'include'})
      .then(resp => {
        if (resp.status !== 200) return rej(`Got bad status code ${resp.status}`);

        res(resp.json());
      })
      .catch(err => rej(err));
  })
}

export const GetMySubmissions = async(hole?: string): Promise<BasicShortSubmission[]|undefined> => {
  return new Promise((res, rej) => {
    let url = `${BackendURL()}/submissions`;
    if (hole) url = `${url}?hole=${hole}`;

    console.log(`*** ${url}`);

    fetch(url, {credentials: 'include'})
      .then(resp => {
        if (resp.status === 401) return res(undefined); // User is not logged in, show banner
        if (resp.status !== 200) return rej(`Got bad status code ${resp.status}`);

        res(resp.json());
      })
      .catch(err => rej(err));
  });
}

// GetMyBestSubmission gets a logged in users best submission using a hole. If the hole is not given, the user isnt logged in,
// or the submission does not exist, undefined will be returned
export const GetMyBestSubmission = async(hole?: string): Promise<BasicShortSubmission|undefined> => {
  return new Promise((res, rej) => {
    if (!hole) return res(undefined);

    const url = `${BackendURL()}/submissions/best/${hole}`;
    console.log(`*** ${url}`);

    fetch(url, {credentials: 'include'})
      .then(resp => {
        if (resp.status === 401 || resp.status === 204) return res(undefined);
        if (resp.status !== 200) return res(undefined);

        res(resp.json());
      })
      // TODO: What is the best way to handle getting the best submission
      .catch(err => rej(`Error getting best submission: ${err}`));
  })
}

export const GetFullSubmission = async(id: string): Promise<BasicFullSubmission|"not logged in"|"not found"> => {
  return new Promise((res, rej) => {
    let url = `${BackendURL()}/submissions/${id}`;
    if (!id) return;

    console.log(`*** ${url}`);

    fetch(url, {credentials: 'include'})
      .then(resp => {
        if (resp.status === 401) return res("not logged in");
        if (resp.status === 404) return res("not found");
        if (resp.status !== 200) return rej(`got bad status code ${resp.status}`);
        return res(resp.json());
      })
      .catch(err => rej(err));
  })
}