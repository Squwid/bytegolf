import { createMuiTheme, ThemeProvider } from '@material-ui/core';
import React from 'react';
import { QueryClient, QueryClientProvider } from 'react-query';
import { PrimaryColor, SecondaryColor } from './Globals';
import Router from './Router';

const theme = createMuiTheme({
  palette: {
    primary: {
      main: PrimaryColor
    },
    secondary: {
      main: SecondaryColor
    }
  }
})

const App: React.FC = () => {
  const qc = new QueryClient({
    defaultOptions: {
      queries: {
        refetchOnWindowFocus: false
      }
    }
  });
  return (
    <QueryClientProvider client={qc}>
      <ThemeProvider theme={theme}>
        <Router />
      </ThemeProvider>
    </QueryClientProvider>
  )
}

export default App;
