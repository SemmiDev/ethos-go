import { useState, useEffect } from 'react';
import apiClient from '../api/client';

const urlBase64ToUint8Array = (base64String) => {
  const padding = '='.repeat((4 - (base64String.length % 4)) % 4);
  const base64 = (base64String + padding).replace(/\-/g, '+').replace(/_/g, '/');

  const rawData = window.atob(base64);
  const outputArray = new Uint8Array(rawData.length);

  for (let i = 0; i < rawData.length; ++i) {
    outputArray[i] = rawData.charCodeAt(i);
  }
  return outputArray;
};

export const usePushNotifications = () => {
  const [isSupported, setIsSupported] = useState(false);
  const [isSubscribed, setIsSubscribed] = useState(false);
  const [permission, setPermission] = useState('default');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if ('serviceWorker' in navigator && 'PushManager' in window) {
      setIsSupported(true);
      setPermission(Notification.permission);

      // Check if already subscribed
      navigator.serviceWorker.ready.then((registration) => {
        registration.pushManager.getSubscription().then((subscription) => {
          setIsSubscribed(!!subscription);
          setLoading(false);
        });
      });
    } else {
      setLoading(false);
    }
  }, []);

  const subscribe = async () => {
    if (!isSupported) return;

    try {
      setLoading(true);

      // 1. Get VAPID Public Key from backend (uses baseURL from apiClient)
      // Response structure: { success: true, data: { vapid_public_key: ... } }
      const response = await apiClient.get('/push/vapid-public-key');
      // apiClient.get returns the axios response object.
      // We want response.data which is the standard envelope.
      const vapidPublicKey = response.data.data.vapid_public_key;

      // 2. Request permission
      const permissionResult = await Notification.requestPermission();
      setPermission(permissionResult);

      if (permissionResult !== 'granted') {
        throw new Error('Permission denied');
      }

      // 3. Register Service Worker (ensure it's registered)
      const registration = await navigator.serviceWorker.register('/sw.js');
      await navigator.serviceWorker.ready;

      // 4. Subscribe to PushManager
      const subscription = await registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: urlBase64ToUint8Array(vapidPublicKey),
      });

      // 5. Send subscription to backend
      // Subscription object needs to be serialized manually or use JSON.stringify then parse
      const subscriptionJSON = subscription.toJSON();

      await apiClient.post('/push/subscribe', {
        endpoint: subscriptionJSON.endpoint,
        keys: {
          p256dh: subscriptionJSON.keys.p256dh,
          auth: subscriptionJSON.keys.auth,
        },
      });

      setIsSubscribed(true);
      console.log('Successfully subscribed to push notifications');
    } catch (error) {
      console.error('Failed to subscribe:', error);
      // If we failed at backend step, we might want to unsubscribe locally to keep state consistent?
      // But maybe not, we can retry later.
    } finally {
      setLoading(false);
    }
  };

  const unsubscribe = async () => {
    if (!isSupported) return;

    try {
      setLoading(true);
      const registration = await navigator.serviceWorker.ready;
      const subscription = await registration.pushManager.getSubscription();

      if (subscription) {
        // Unsubscribe from backend first
        await apiClient.post('/push/unsubscribe', {
          endpoint: subscription.endpoint,
        });

        // Unsubscribe from browser
        await subscription.unsubscribe();
        setIsSubscribed(false);
      }
    } catch (error) {
      console.error('Failed to unsubscribe:', error);
    } finally {
      setLoading(false);
    }
  };

  return {
    isSupported,
    isSubscribed,
    permission,
    loading,
    subscribe,
    unsubscribe,
  };
};
