# Ethos - Habit Tracker Frontend

A beautiful, modern React frontend for the Ethos Habit Tracking application. Built with Vite, React Router, Tailwind CSS, and Zustand for state management.

![Ethos Screenshot](./screenshot.png)

## ✨ Features

-   🔐 **Authentication** - Register, login, logout, and session management
-   📊 **Dashboard** - Overview of your habits with quick actions
-   🎯 **Habit Management** - Create, edit, delete, activate/deactivate habits
-   📝 **Activity Logging** - Log your habit completions with notes
-   📈 **Analytics** - Track your progress with visual charts and achievements
-   ⚙️ **Settings** - Manage your profile and preferences
-   🌙 **Dark Theme** - Beautiful glassmorphism design with purple gradients
-   📱 **Responsive** - Works on desktop, tablet, and mobile

## 🛠️ Tech Stack

-   **React 18** - UI Library
-   **Vite** - Build tool and dev server
-   **React Router v6** - Client-side routing
-   **Zustand** - State management
-   **Tailwind CSS v4** - Utility-first CSS framework
-   **Axios** - HTTP client
-   **Headless UI** - Unstyled accessible components
-   **Lucide React** - Beautiful icons
-   **date-fns** - Date manipulation

## 🚀 Getting Started

### Prerequisites

-   Node.js 18+
-   npm or yarn
-   Backend API running on `http://localhost:8080`

### Installation

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

The app will be available at `http://localhost:3000`

### Build for Production

```bash
npm run build
```

The build output will be in the `dist` folder.

## 📁 Project Structure

```
frontend/
├── src/
│   ├── api/                    # API client and service modules
│   │   ├── client.js           # Axios instance with interceptors
│   │   ├── auth.js             # Auth API endpoints
│   │   └── habits.js           # Habits API endpoints
│   │
│   ├── components/             # Reusable components
│   │   ├── ui/                 # Base UI components
│   │   │   ├── Button.jsx
│   │   │   ├── Input.jsx
│   │   │   ├── Modal.jsx
│   │   │   ├── Card.jsx
│   │   │   ├── Badge.jsx
│   │   │   ├── Toast.jsx
│   │   │   └── Loading.jsx
│   │   │
│   │   ├── layout/             # Layout components
│   │   │   ├── Sidebar.jsx
│   │   │   └── Layout.jsx
│   │   │
│   │   └── habits/             # Habit-specific components
│   │       ├── HabitCard.jsx
│   │       ├── HabitModals.jsx
│   │       └── HabitLogList.jsx
│   │
│   ├── pages/                  # Page components
│   │   ├── auth/
│   │   │   └── AuthPages.jsx   # Login & Register
│   │   ├── dashboard/
│   │   │   └── DashboardPage.jsx
│   │   ├── habits/
│   │   │   ├── HabitsListPage.jsx
│   │   │   └── HabitDetailPage.jsx
│   │   ├── analytics/
│   │   │   └── AnalyticsPage.jsx
│   │   └── settings/
│   │       └── SettingsPage.jsx
│   │
│   ├── stores/                 # Zustand state stores
│   │   ├── authStore.js        # Authentication state
│   │   ├── habitsStore.js      # Habits state & actions
│   │   └── uiStore.js          # UI state (toasts, modals)
│   │
│   ├── App.jsx                 # Main app with routing
│   ├── main.jsx                # Entry point
│   └── index.css               # Global styles & design system
│
├── index.html
├── vite.config.js
├── package.json
└── README.md
```

## 🎨 Design System

The app uses a custom design system with:

-   **Glassmorphism** - Frosted glass effect with blur
-   **Gradient Accents** - Purple to teal gradient theme
-   **Dark Mode** - Eye-friendly dark background
-   **Micro-animations** - Smooth transitions and hover effects
-   **Responsive** - Mobile-first responsive design

### Color Palette

| Color       | Hex       | Usage             |
| ----------- | --------- | ----------------- |
| Violet 500  | `#8b5cf6` | Primary actions   |
| Teal 400    | `#22d3d1` | Accent color      |
| Emerald 500 | `#22c55e` | Success states    |
| Amber 500   | `#f59e0b` | Warnings, streaks |
| Red 400     | `#f87171` | Errors, danger    |

## 📡 API Integration

The frontend connects to the backend API with:

-   **Base URL**: `/api` (proxied to `http://localhost:8080`)
-   **Authentication**: JWT Bearer tokens
-   **Auto-refresh**: Tokens are stored in localStorage

### Available Endpoints

| Endpoint                | Description            |
| ----------------------- | ---------------------- |
| `POST /auth/register`   | Register new user      |
| `POST /auth/login`      | Login user             |
| `POST /auth/logout`     | Logout current session |
| `POST /auth/logout-all` | Logout all devices     |
| `GET /habits`           | List all habits        |
| `POST /habits`          | Create new habit       |
| `GET /habits/:id`       | Get habit details      |
| `PUT /habits/:id`       | Update habit           |
| `DELETE /habits/:id`    | Delete habit           |
| `POST /habits/:id/logs` | Log habit completion   |
| `GET /habits/:id/stats` | Get habit statistics   |
| `GET /dashboard`        | Get dashboard data     |

## 🔧 Configuration

### Environment Variables

Create a `.env` file in the frontend directory:

```env
VITE_API_URL=http://localhost:8080/api
```

### Proxy Configuration

The Vite dev server proxies API requests to the backend:

```js
// vite.config.js
export default defineConfig({
    server: {
        proxy: {
            '/api': {
                target: 'http://localhost:8080',
                changeOrigin: true,
            },
        },
    },
});
```

## 📱 Pages Overview

### Login/Register

-   Email and password authentication
-   Form validation
-   Beautiful animated UI

### Dashboard

-   Stats cards showing active habits, completions, and streaks
-   Today's habits with quick log actions
-   Motivational tips

### Habits List

-   Grid/List view toggle
-   Filter by status (all/active/inactive)
-   Search functionality
-   Quick actions menu

### Habit Detail

-   Full habit information
-   Statistics (total logs, current/longest streak)
-   Activity log history
-   Edit and delete options

### Analytics

-   Weekly progress chart
-   Habit distribution by frequency
-   Achievements system

### Settings

-   Profile information
-   Notification preferences
-   Security options (logout, logout all)

## 🧪 Testing

```bash
# Run linting
npm run lint
```

## 📄 License

MIT License

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

---

Built with ❤️ using React and Vite
