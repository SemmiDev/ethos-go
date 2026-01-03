self.addEventListener('push', function (event) {
  if (event.data) {
    const payload = event.data.json();

    const options = {
      body: payload.body,
      icon: payload.icon || '/icon-192x192.png',
      badge: payload.badge || '/badge-72x72.png',
      data: payload.data,
      actions: payload.actions,
      tag: payload.tag,
      renotify: !!payload.tag,
    };

    event.waitUntil(self.registration.showNotification(payload.title, options));
  }
});

self.addEventListener('notificationclick', function (event) {
  event.notification.close();

  // Handle action clicks
  if (event.action) {
    // You can handle specific actions here
    console.log('Action clicked:', event.action);
  }

  // Open the app when notification is clicked
  event.waitUntil(
    clients
      .matchAll({
        type: 'window',
        includeUncontrolled: true,
      })
      .then(function (clientList) {
        // Check if there's already a window open
        for (let i = 0; i < clientList.length; i++) {
          const client = clientList[i];
          if (client.url.includes('/') && 'focus' in client) {
            return client.focus();
          }
        }
        // If no window is open, open a new one
        if (clients.openWindow) {
          return clients.openWindow('/');
        }
      })
  );
});
