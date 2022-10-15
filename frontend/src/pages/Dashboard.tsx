import { Delete, Memory, TableRows } from '@mui/icons-material';
import {
  AppBar,
  Avatar,
  Container,
  IconButton,
  List,
  ListItem,
  ListItemAvatar,
  ListItemText,
  Toolbar,
} from '@mui/material';
import { Fragment, useEffect, useState } from 'react';
import { api } from '../api';
import * as apiGen from '../apiCodegen';
import { HeadLabel } from './components/HeadLabel';

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
          <HeadLabel>Dashboard</HeadLabel>
        </Toolbar>
      </AppBar>
      {err ? (
        <div>Error: ${err.message}</div>
      ) : (
        <Container maxWidth="md">
          <List>
            {tables.map((e) => {
              return (
                <TableItem
                  tableName={e.tableName}
                  updateTables={updateTables}
                />
              );
            })}
          </List>
        </Container>
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
  return (
    <ListItem
      secondaryAction={
        <IconButton
          onClick={() => {
            api.dropTable(tableName);
            updateTables();
          }}
          edge="end"
          aria-label="delete"
        >
          <Delete />
        </IconButton>
      }
    >
      <ListItemAvatar>
        <Avatar>
          <TableRows />
        </Avatar>
      </ListItemAvatar>
      <ListItemText primary={tableName} />
    </ListItem>
  );
};
