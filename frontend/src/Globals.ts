export const PrimaryColor = '#0baa54';
export const SecondaryColor = '#de5353';
export const ThirdColor = '#414141';
export const LightTextColor = '#e6e6e6';

const backend = process.env.REACT_APP_BACKEND_URI ? process.env.REACT_APP_BACKEND_URI : 'http://192.168.0.21:9999';
export const readme = process.env.REACT_APP_README_URI ? process.env.REACT_APP_README_URI : 'https://raw.githubusercontent.com/Squwid/bytegolf/master/README.md';


export const BackendURL = () => `${backend}/api`
export const RawBackendURL = () => `${backend}`;