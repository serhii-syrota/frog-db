import {
  Container,
  CssBaseline,
  Box,
  TextField,
  Button,
  Avatar,
} from '@mui/material';
import { redirect } from 'react-router-dom';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';

export const SetApiUrl = () => {
  const apiUrlId = 'apiUrl';
  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const data = new FormData(event.currentTarget);
    const apiUrl = data.get(apiUrlId);
    redirect('dashboard');
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
