#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

use tauri::Manager;

mod token;

/// Starts the Courrier window with sign-in deep links, notifications, the system browser, self-update and the keychain-backed session token wired in.
fn main() {
    tauri::Builder::default()
        .plugin(tauri_plugin_single_instance::init(|app, _argv, _cwd| {
            if let Some(window) = app.get_webview_window("main") {
                let _ = window.unminimize();
                let _ = window.set_focus();
            }
        }))
        .plugin(tauri_plugin_deep_link::init())
        .plugin(tauri_plugin_notification::init())
        .plugin(tauri_plugin_opener::init())
        .plugin(tauri_plugin_updater::Builder::new().build())
        .plugin(tauri_plugin_process::init())
        .invoke_handler(tauri::generate_handler![
            token::token_get,
            token::token_set,
            token::token_clear
        ])
        .setup(|_app| {
            #[cfg(any(target_os = "linux", windows))]
            {
                use tauri_plugin_deep_link::DeepLinkExt;
                if let Err(err) = _app.deep_link().register_all() {
                    eprintln!("courrier: could not register the sign-in link handler: {err}");
                }
            }
            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("Courrier failed to start");
}
