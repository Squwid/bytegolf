export type Language = {
  name: string;
  editorValue: string;
  language: string;
  version: string;
}

export const getLanguage = (lang: string): Language => {
  for (const l of AvailableLanguages) {
    if (l.language === lang) return l;
  }
  return AvailableLanguages[0];
}

export const AvailableLanguages: Language[] = [
  {
    name: 'Python 3',
    editorValue: 'python',
    language: 'python3',
    version: '3.12.2'
  },
  {
    name: 'PHP',
    editorValue: 'php',
    language: 'php',
    version: '7.3.10'
  },
  {
    name: 'Javascript (node 22)',
    editorValue: 'javascript',
    language: 'node',
    version: '22'
  },
  {
    name: 'Go',
    editorValue: 'golang',
    language: 'go',
    version: '1.22.2'
  },
  {
    name: 'Bash',
    editorValue: 'batchfile',
    language: 'bash',
    version: '5.2.26'
  }


  // {
  //   name: 'Java',
  //   editorValue: 'java',
  //   language: 'java',
  //   version: 'JDK 11.0.4'
  // },
  // {
  //   name: 'C++',
  //   editorValue: 'c_cpp',
  //   language: 'c++',
  //   version: 'g++ 17 GCC 9.10'
  // },
  // {
  //   name: 'Rust',
  //   editorValue: 'rust',
  //   language: 'rust',
  //   version: '1.38.0'
  // },
  // {
  //   name: 'Typescript',
  //   editorValue: 'typescript',
  //   language: 'typescript'
  // },
  // {
  //   name: 'Powershell',
  //   editorValue: 'powershell',
  //   language: 'powershell'
  // },
]