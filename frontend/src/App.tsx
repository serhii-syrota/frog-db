import '@fontsource/roboto/300.css';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/500.css';
import '@fontsource/roboto/700.css';
import {
  Avatar,
  Box,
  Button,
  Container,
  createTheme,
  TextField,
  ThemeProvider,
} from '@mui/material';

import { Helmet } from 'react-helmet';
import { CssBaseline } from '@mui/material';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import { api } from './api';

const theme = createTheme();
export const App = () => {
  return (
    <div className="App">
      <ThemeProvider theme={theme}>
        <Helmet>
          <meta
            name="viewport"
            content="initial-scale=1, width=device-width"
          />
        </Helmet>
        <CssBaseline />
        <SetSourcePage />
      </ThemeProvider>
    </div>
  );
};

// Enter api url page: fill api link
// Db schema page: buttons drop, add, table view, JSON dump, change url
// Table view page: delete where, select where, insert, update
// Header with home page

export const SetSourcePage = () => {
  api.dbSchema().then(console.log);
  return <SignIn></SignIn>;
};

const SignIn = () => {
  const apiUrlId = 'apiUrl';
  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const data = new FormData(event.currentTarget);
    const apiUrl = data.get(apiUrlId);
    console.log(apiUrl);
  };

  return (
    <Container>
      <CssBaseline />
      <Box
        sx={{
          marginTop: 25,
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
        }}
      >
        <Avatar sx={{ m: 1, bgcolor: 'secondary.main' }}>
          <LockOutlinedIcon />
        </Avatar>
        <Box
          component="form"
          onSubmit={handleSubmit}
          noValidate
          sx={{ mt: 1 }}
        >
          <TextField
            margin="normal"
            required
            fullWidth
            id={apiUrlId}
            label="Frogdb api url"
            name={apiUrlId}
            autoComplete="url"
            autoFocus
          />
          <Button
            type="submit"
            fullWidth
            variant="contained"
            sx={{ mt: 3, mb: 2 }}
          >
            Run dashboard
          </Button>
        </Box>
      </Box>
    </Container>
  );
};
