const API_BASE = 'https://orhalevanah-production.up.railway.app/';

export interface CalendarDay {
  civil_date: string;
  biblical_date?: string;
  biblical_day?: number;
  biblical_month?: number;
  is_today?: boolean;
  is_shabbat?: boolean;
  is_new_month?: boolean;
  feasts?: string[];
}

export interface CalendarResponse {
  success: boolean;
  message?: string;
  data: {
    year: number;
    month: number;
    month_name: string;
    days: CalendarDay[];
  };
}

export interface BiblicalDateResponse {
  success: boolean;
  message?: string;
  data: {
    civil_date: string;
    biblical_date: string;
    jerusalem_time?: string;
    sunset_time?: string;
    after_sunset?: boolean;
    day_note?: string;
    biblical_month?: number;
    biblical_day?: number;
    month_name?: string;
  };
}

export interface FeastsResponse {
  success: boolean;
  message?: string;
  data: {
    current_feasts?: Array<{
      name: string;
      biblical_month: number;
      biblical_day: number;
      gregorian_date?: string;
      description?: string;
    }>;
    upcoming_feasts?: Array<{
      name: string;
      biblical_month: number;
      biblical_day: number;
      gregorian_date?: string;
      description?: string;
    }>;
  };
}

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const text = await response.text();
    throw new Error(`HTTP ${response.status}: ${text}`);
  }
  return response.json();
}

export async function getCalendar(year: number, month: number): Promise<CalendarResponse> {
  const res = await fetch(`${API_BASE}/api/calendar?year=${year}&month=${month}`);
  return handleResponse<CalendarResponse>(res);
}

export async function getBiblicalDate(date: string): Promise<BiblicalDateResponse> {
  const res = await fetch(`${API_BASE}/api/biblical-date?date=${date}`);
  return handleResponse<BiblicalDateResponse>(res);
}

export async function getFeasts(date: string): Promise<FeastsResponse> {
  const res = await fetch(`${API_BASE}/api/feasts?date=${date}`);
  return handleResponse<FeastsResponse>(res);
}