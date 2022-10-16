import { Delete, Memory, TableRows } from '@mui/icons-material';
import {
  AppBar,
  Avatar,
  Button,
  Container,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  IconButton,
  List,
  ListItem,
  ListItemAvatar,
  ListItemButton,
  ListItemText,
  Toolbar,
  Tooltip,
} from '@mui/material';
import { Fragment, useEffect, useState } from 'react';
import { api } from '../api';
import * as apiGen from '../apiCodegen';
import { AddTable } from './CreateTable';
import { useNavigate } from 'react-router-dom';
import { HeadLabelClickable } from './components/HeadLabel';

export const Dashboard = () => {
  const [tables, setTables] = useState(
    [] as Required<apiGen.TableSchema>[]
  );
  const [err, setErr] = useState<Error | null>(null);
  const updateTables = () => {
    api.dbSchema().then((res) => {
      if (res[0]) {
        return setErr(res[0]);
      }
      return setTables(res[1]);
    });
  };
  useEffect(() => {
    updateTables();
  }, []);
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
      {err ? (
        <div>Error: ${err.message}</div>
      ) : (
        <>
          <Container maxWidth="md">
            <List>
              {tables
                .sort((a, b) =>
                  a.tableName.localeCompare(b.tableName)
                )
                .map((e, i) => {
                  return (
                    <TableItem
                      key={i}
                      tableName={e.tableName}
                      updateTables={updateTables}
                    />
                  );
                })}
            </List>
          </Container>
          <AddTable updateTables={updateTables} />
        </>
      )}
    </Fragment>
  );
};

const TableItem = ({
  tableName,
  updateTables,
}: {
  tableName: string;
  updateTables: () => void;
}) => {
  const navigate = useNavigate();
  const toTablePage = () => {
    navigate(`/table/${tableName}`);
  };

  return (
    <ListItem
      secondaryAction={
        <DeleteTable
          tableName={tableName}
          dropTable={() => {
            api.dropTable(tableName);
            updateTables();
          }}
        />
      }
    >
      <ListItemButton
        onClick={() => {
          toTablePage();
        }}
      >
        <ListItemAvatar>
          <Avatar>
            <TableRows />
          </Avatar>
        </ListItemAvatar>
        <ListItemText primary={tableName} />
      </ListItemButton>
    </ListItem>
  );
};

const DeleteTable = ({
  tableName,
  dropTable,
}: {
  tableName: string;
  dropTable: () => void;
}) => {
  const [open, setOpen] = useState(false);
  const handleClickOpen = () => {
    setOpen(true);
  };
  const handleClose = () => {
    setOpen(false);
  };

  return (
    <div>
      <Tooltip title="Delete">
        <IconButton
          onClick={handleClickOpen}
          edge="end"
          aria-label="delete"
        >
          <Delete />
        </IconButton>
      </Tooltip>

      <Dialog
        open={open}
        onClose={handleClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
        <DialogTitle id="alert-dialog-title">
          {`Are you sure to drop table ${tableName}?`}
        </DialogTitle>
        <DialogContent>
          <DialogContentText id="alert-dialog-description">
            You might be able to restore the table after deletion.
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose}>Cancel</Button>
          <Button
            onClick={() => {
              handleClose();
              dropTable();
            }}
            color={'error'}
            autoFocus
          >
            Drop
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};
