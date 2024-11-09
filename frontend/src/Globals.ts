export const PrimaryColor = "#0baa54";
export const SecondaryColor = "#de5353";
export const ThirdColor = "#414141";
export const LightTextColor = "#e6e6e6";

const backend = "https://api.play.byte.golf";
export const readme = process.env.REACT_APP_README_URI
  ? process.env.REACT_APP_README_URI
  : "https://raw.githubusercontent.com/Squwid/bytegolf/master/README.md";

export const BackendURL = () => `${backend}/api`;
export const RawBackendURL = () => `${backend}`;
