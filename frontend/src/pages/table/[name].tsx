import { Memory } from '@mui/icons-material';
import {
  AppBar,
  Button,
  ImageListItem,
  Toolbar,
} from '@mui/material';
import { AxiosError } from 'axios';
import { Fragment, useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { api } from '../../api';
import { Info } from '../../apiCodegen';
import { HeadLabelClickable } from '../components/HeadLabel';
import * as apiGen from '../../apiCodegen';
import { DataGrid, GridColDef } from '@mui/x-data-grid';

export const TableView = () => {
  const [tableSchema, setTableSchema] = useState<apiGen.Schema[]>([]);
  const [rows, setRows] = useState<Record<string, unknown>[]>([]);
  const { name } = useParams<{ name: string }>();

  const deleteDuplicatesClick = () => {
    removeDuplicates(name!);
    uploadRows(setRows, name!);
  };

  useEffect(() => {
    uploadRows(setRows, name!);
    uploadSchema(setTableSchema, name!);
  }, [name]);
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

      <div style={{ height: 630, width: '100%' }}>
        <DataGrid
          rows={rows.map((e, i) => ({ ...e, id: i }))}
          columns={schemaToColDef(tableSchema)}
          autoPageSize
          rowsPerPageOptions={[10]}
          density="comfortable"
        />
        <Button
          sx={{ bottom: '3%', right: '44%', position: 'fixed' }}
          variant="contained"
          onClick={deleteDuplicatesClick}
        >
          Remove duplicates
        </Button>
      </div>
    </Fragment>
  );
};

const uploadRows = (
  setRows: (val: Record<string, unknown>[]) => void,
  name: string
): void => {
  api
    .selectRows(name!, { columns: [], conditions: {} })
    .then((res) => {
      if (res[0]) {
        alert(
          ((res[0] as AxiosError).response!.data as Info).message
        );
        return;
      }
      setRows(res[1]);
    });
};
const uploadSchema = (
  setTableSchema: (val: apiGen.Schema[]) => void,
  name: string
) => {
  api.dbSchema().then((res) => {
    if (res[0]) {
      alert(((res[0] as AxiosError).response!.data as Info).message);
      return;
    }
    setTableSchema(res[1].find((e) => e.tableName === name!)!.schema);
  });
};

const schemaToColDef = (schema: apiGen.Schema[]): GridColDef[] => {
  return [
    ...schema.map(
      (e): GridColDef => ({
        field: e.column,
        headerName: e.column + ' (' + e.type + ')',
        width: 150,
        type: typeToColDef(e.type),
        renderCell: (params) => {
          if (e.type === apiGen.SchemaTypeEnum.Image) {
            return (
              <ImageListItem>
                <img src={params.value} loading="lazy" alt="cell" />
              </ImageListItem>
            );
          }
          if (e.type === apiGen.SchemaTypeEnum.RealInv) {
            return (
              <div>
                ({params.value[0]}, {params.value[1]})
              </div>
            );
          }
          return <div>{params.value}</div>;
        },
        filterable: e.type !== apiGen.SchemaTypeEnum.RealInv,
        sortable:
          e.type !== apiGen.SchemaTypeEnum.Image &&
          e.type !== apiGen.SchemaTypeEnum.RealInv,
      })
    ),
  ];
};

const typeToColDef = (t: string) => {
  if (
    t === apiGen.SchemaTypeEnum.Integer ||
    t === apiGen.SchemaTypeEnum.Real
  ) {
    return 'number';
  }
  return 'string';
};

const removeDuplicates = (name: string) => {
  return api.deleteDuplicateRows(name);
};
