export { request } from './request';
export {
  {{~ for table in tables ~}}
  {{ table.label|dbcore_capitalize }},
  {{~ end ~}}
} from './types';
