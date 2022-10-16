import '@fontsource/roboto/300.css';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/500.css';
import '@fontsource/roboto/700.css';
import { createTheme, ThemeProvider } from '@mui/material';
import {
  createBrowserRouter,
  RouterProvider,
} from 'react-router-dom';
import { Helmet } from 'react-helmet';
import { CssBaseline } from '@mui/material';
import { SetApiUrl } from './pages/SetApiUrl';
import { Dashboard } from './pages/Dashboard';
import { TableView } from './pages/table/[name]';

// Enter api url page: fill api link
// Db schema page: buttons drop, add, table view, JSON dump, change url
// Table view page: delete where, select where, insert, update
// Header with home page

const router = createBrowserRouter([
  {
    path: '/',
    element: <SetApiUrl />,
  },
  {
    path: '/dashboard',
    element: <Dashboard />,
  },
  {
    path: '/table/:name',
    element: <TableView />,
  },
]);

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
        <RouterProvider router={router} />
      </ThemeProvider>
    </div>
  );
};
