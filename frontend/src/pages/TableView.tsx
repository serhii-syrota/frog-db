import { Memory } from '@mui/icons-material';
import { AppBar, Toolbar } from '@mui/material';
import { Fragment } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { HeadLabelClickable } from './components/HeadLabel';
export const TableView = () => {
  const navigate = useNavigate();
  const toDashboard = () => {
    navigate(`/dashboard`);
  };

  const { name } = useParams<{ name: string }>();
  if (!name) {
    toDashboard();
    return <></>;
  }
  return (
    <Fragment>
      <AppBar position="relative">
        <Toolbar>
          <Memory />
          <HeadLabelClickable href="/dashboard">
            Dashboard
          </HeadLabelClickable>
        </Toolbar>
      </AppBar>
    </Fragment>
  );
};
