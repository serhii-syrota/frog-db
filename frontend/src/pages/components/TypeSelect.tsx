import { NativeSelect } from '@mui/material';
import { BootstrapInput } from './BootstrapInput';
import * as apiGen from '../../apiCodegen';

export const TypeSelect = ({
  onChange,
}: {
  onChange: (v: apiGen.SchemaTypeEnum) => void;
}) => {
  return (
    <NativeSelect
      sx={{ mt: 1 }}
      onChange={(e) => {
        onChange(e.target.value as apiGen.SchemaTypeEnum);
      }}
      input={<BootstrapInput />}
    >
      {Object.keys(apiGen.SchemaTypeEnum).map((key) => (
        <option
          value={
            apiGen.SchemaTypeEnum[
              key as keyof typeof apiGen.SchemaTypeEnum
            ]
          }
        >
          {key}
        </option>
      ))}
    </NativeSelect>
  );
};
