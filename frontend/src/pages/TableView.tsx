import { useLoaderData } from 'react-router-dom';
export const TableView = () => {
  const data = useLoaderData();
  console.log(data);

  return <div></div>;
};
