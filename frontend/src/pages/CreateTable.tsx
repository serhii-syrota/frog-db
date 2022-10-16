import { Add, Close } from '@mui/icons-material';
import {
  Slide,
  Fab,
  Dialog,
  DialogContent,
  AppBar,
  Button,
  IconButton,
  Toolbar,
  Typography,
  FormControl,
} from '@mui/material';
import { TransitionProps } from '@mui/material/transitions';
import React from 'react';
import { forwardRef, useState } from 'react';
import { BootstrapInput } from './components/BootstrapInput';
import { GridCenter } from './components/GridCenter';
import { TypeSelect } from './components/TypeSelect';
import * as apiGen from '../apiCodegen';
import { api } from '../api';
import { AxiosError } from 'axios';

const Transition = forwardRef(
  (
    props: TransitionProps & {
      children: React.ReactElement;
    },
    ref: React.Ref<unknown>
  ) => {
    return <Slide direction="up" ref={ref} {...props} />;
  }
);

export const AddTable = ({
  updateTables,
}: {
  updateTables: () => void;
}) => {
  const [newTable, setNewTable] = useState<
    Required<apiGen.TableSchema>
  >(tableSchemaDefaultVal());
  const [open, setOpen] = useState(false);
  const handleClickOpen = () => {
    setOpen(true);
  };
  const handleClose = () => {
    setOpen(false);
  };
  return (
    <div>
      <AddTableFab handleClickOpen={handleClickOpen} />
      <Dialog
        open={open}
        onClose={handleClose}
        fullScreen
        TransitionComponent={Transition}
      >
        <HeaderBar
          handleClose={handleClose}
          handleSave={() => {
            const validationErr = validateNewTable(newTable);
            if (validationErr) {
              alert(validationErr.message);
              return;
            }
            api.createTable(newTable).then((res) => {
              if (res[0]) {
                alert(
                  (
                    (res[0] as AxiosError).response!.data as {
                      message: string;
                    }
                  ).message
                );
                return;
              }
              setNewTable(tableSchemaDefaultVal());
              updateTables();
              handleClose();
            });
          }}
        />
        <DialogContent>
          <AddForm newTable={newTable} setNewTable={setNewTable} />
        </DialogContent>
      </Dialog>
    </div>
  );
};

const AddForm = ({
  newTable,
  setNewTable,
}: {
  newTable: Required<apiGen.TableSchema>;
  setNewTable: (table: Required<apiGen.TableSchema>) => void;
}) => {
  return (
    <>
      <GridCenter>
        <FormControl sx={{ m: 1 }} variant="standard">
          <Typography variant="h6" align="center">
            Table name
          </Typography>
          <BootstrapInput
            onChange={(e) => {
              setNewTable({ ...newTable, tableName: e.target.value });
            }}
            sx={{ mt: 1 }}
          />
        </FormControl>
      </GridCenter>
      <GridCenter>
        <FormControl sx={{ m: 1 }} variant="standard">
          <Typography variant="h6" align="center">
            Column name
          </Typography>
          {newTable.schema.map((_, i) => {
            return (
              <BootstrapInput
                onChange={(v) => {
                  newTable['schema'][i].column = v.target.value;
                  setNewTable(newTable);
                }}
                key={i}
                sx={{ mt: 1 }}
              />
            );
          })}
        </FormControl>
        <FormControl sx={{ m: 1 }} variant="standard">
          <Typography variant="h6" align="center">
            Type
          </Typography>
          {newTable.schema.map((_, i) => {
            return (
              <TypeSelect
                onChange={(v) => {
                  newTable['schema'][i].type = v;
                  setNewTable(newTable);
                }}
                key={i}
              />
            );
          })}
        </FormControl>
      </GridCenter>
      <GridCenter>
        <Fab
          onClick={() => {
            setNewTable({
              tableName: newTable.tableName,
              schema: [
                ...newTable.schema,
                { column: '', type: 'integer' },
              ],
            });
          }}
          color="primary"
        >
          <Add />
        </Fab>
      </GridCenter>
    </>
  );
};

const HeaderBar = ({
  handleClose,
  handleSave,
}: {
  handleClose: () => void;
  handleSave: () => void;
}) => {
  return (
    <AppBar sx={{ position: 'relative' }}>
      <Toolbar>
        <HeaderClose handleClose={handleClose} />
        <Button
          color="inherit"
          size="large"
          onClick={handleSave}
          component="div"
        >
          Save
        </Button>
      </Toolbar>
    </AppBar>
  );
};

const HeaderClose = ({
  handleClose,
}: {
  handleClose: () => void;
}) => {
  return (
    <>
      <IconButton
        edge="start"
        color="inherit"
        onClick={handleClose}
        aria-label="close"
      >
        <Close />
      </IconButton>
      <Typography
        sx={{ ml: 2, flex: 1 }}
        variant="h6"
        component="div"
      >
        Create table
      </Typography>
    </>
  );
};

const AddTableFab = ({
  handleClickOpen,
}: {
  handleClickOpen: () => void;
}) => {
  return (
    <Fab
      onClick={handleClickOpen}
      sx={{ position: 'fixed', bottom: 16, right: 16 }}
      color="primary"
      aria-label="add"
      variant="extended"
    >
      <Add />
      Add table
    </Fab>
  );
};

const validateNewTable = (
  table: Required<apiGen.TableSchema>
): Error | null => {
  if (table.tableName.length < 5)
    return new Error('Table name length min 5 symbols');

  return null;
};

const tableSchemaDefaultVal = () => ({
  tableName: '',
  schema: [{ column: '', type: apiGen.SchemaTypeEnum.Integer }],
});
